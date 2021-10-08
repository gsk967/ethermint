package testutil

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tharsis/ethermint/x/nameservice/client/cli"
	"github.com/tharsis/ethermint/x/nameservice/types"
	"os"
	"time"
)

func (s *IntegrationTestSuite) TestGetCmdQueryParams() {
	val := s.network.Validators[0]
	sr := s.Require()

	testCases := []struct {
		name string
		args []string
	}{
		{
			"params",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := cli.GetQueryParamsCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			sr.NoError(err)
			var param types.QueryParamsResponse
			err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &param)
			sr.NoError(err)
			params := types.DefaultParams()
			params.RecordRent = sdk.NewCoin(s.cfg.BondDenom, types.DefaultRecordRent)
			params.RecordRentDuration = 5 * time.Second
			params.AuthorityGracePeriod = 5 * time.Second
			sr.Equal(param.GetParams().String(), params.String())
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryForRecords() {
	val := s.network.Validators[0]
	sr := s.Require()
	var recordID string
	var bondId string

	testCases := []struct {
		name        string
		args        []string
		expErr      bool
		noOfRecords int
		preRun      func()
	}{
		{
			"invalid request",
			[]string{"invalid", fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			true,
			0,
			func() {

			},
		},
		{
			"get records list",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			false,
			1,
			func() {
				CreateBond(s)
				// get the bond id from bond list
				bondId := GetBondId(s)
				dir, err := os.Getwd()
				sr.NoError(err)
				payloadPath := dir + "/example1.yml"
				args := []string{
					payloadPath, bondId,
					fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
					fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
					fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
					fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
				}

				clientCtx := val.ClientCtx
				cmd := cli.GetCmdSetRecord()

				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			clientCtx := val.ClientCtx
			if !tc.expErr {
				tc.preRun()
			}
			cmd := cli.GetCmdList()
			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var records []types.RecordType
				err := json.Unmarshal(out.Bytes(), &records)
				sr.NoError(err)
				sr.Equal(tc.noOfRecords, len(records))
				recordID = records[0].Id
				bondId = GetBondId(s)
			}
		})
	}

	s.T().Log("Test Cases for getting records by record-id")
	testCasesByRecordID := []struct {
		name   string
		args   []string
		expErr bool
	}{
		{
			"invalid request without record id",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			true,
		},
		{
			"get records by id",
			[]string{recordID, fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			false,
		},
	}

	for _, tc := range testCasesByRecordID {
		s.Run(tc.name, func() {
			cmd := cli.GetCmdGetResource()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var response types.QueryRecordByIdResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &response)
				sr.NoError(err)
				sr.NotNil(response.GetRecord())
			}
		})
	}

	s.T().Log("Test Cases for getting records by bond-id")
	testCasesByRecordByBondID := []struct {
		name   string
		args   []string
		expErr bool
	}{
		{
			"invalid request without bond-id",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			true,
		},
		{
			"get records by bond-id",
			[]string{bondId, fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			false,
		},
	}

	for _, tc := range testCasesByRecordByBondID {
		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryByBond()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var response types.QueryRecordByBondIdResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &response)
				sr.NoError(err)
			}
		})
	}

	s.T().Log("Test Cases for getting nameservice module account balance")
	testCasesForNameServiceModuleBalance := []struct {
		name        string
		args        []string
		expErr      bool
		noOfRecords int
	}{
		{
			"get nameservice module accounts balance",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			false,
			1,
		},
	}

	for _, tc := range testCasesForNameServiceModuleBalance {
		s.Run(tc.name, func() {
			cmd := cli.GetCmdBalance()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var response types.GetNameServiceModuleBalanceResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &response)
				sr.NoError(err)
				sr.Equal(tc.noOfRecords, len(response.GetBalances()))
				balance := response.GetBalances()[0]
				sr.Equal(balance.AccountName, types.RecordRentModuleAccountName)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdWhoIs() {
	val := s.network.Validators[0]
	sr := s.Require()
	var authorityName = "test2"
	testCases := []struct {
		name        string
		args        []string
		expErr      bool
		noOfRecords int
		preRun      func(authorityName string)
	}{
		{
			"invalid request without name",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			true,
			1,
			func(authorityName string) {

			},
		},
		{
			"success query with name",
			[]string{authorityName, fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			false,
			1,
			func(authorityName string) {
				// reserving the name
				clientCtx := val.ClientCtx
				cmd := cli.GetCmdReserveName()
				args := []string{
					authorityName,
					fmt.Sprintf("--owner=%s", accountAddress),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
					fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
					fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
					fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
				}
				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if !tc.expErr {
				tc.preRun(authorityName)
			}
			cmd := cli.GetCmdWhoIs()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var response types.QueryWhoisResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &response)
				sr.NoError(err)
				nameAuthority := response.GetNameAuthority()
				nameAuthority.OwnerAddress = accountAddress
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdLookupWRN() {
	val := s.network.Validators[0]
	sr := s.Require()
	var authorityName = "test1"
	testCases := []struct {
		name        string
		args        []string
		expErr      bool
		noOfRecords int
		preRun      func(authorityName string)
	}{
		{
			"invalid request without wrn",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			true,
			0,
			func(authorityName string) {

			},
		},
		{
			"success query with name",
			[]string{fmt.Sprintf("wrn://%s/", authorityName), fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			false,
			1,
			func(authorityName string) {
				// reserving the name
				createNameRecord(authorityName, s)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if !tc.expErr {
				// set-name with wrn and bond-id
				tc.preRun(authorityName)
			}
			cmd := cli.GetCmdLookupWRN()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var response types.QueryLookupWrnResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &response)
				sr.NoError(err)
				nameRecord := response.GetName()
				sr.NotNil(nameRecord.Latest.Id)
			}
		})
	}

	testCasesForNamesList := []struct {
		name        string
		args        []string
		expErr      bool
		noOfRecords int
	}{
		{
			"invalid request without wrn",
			[]string{"invalid", fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			true,
			0,
		},
		{
			"success query with name",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			false,
			1,
		},
	}

	for _, tc := range testCasesForNamesList {
		s.Run(tc.name, func() {
			cmd := cli.GetCmdNames()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var response types.QueryListNameRecordsResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &response)
				sr.NoError(err)
				sr.NotZero(len(response.GetNames()))
			}
		})
	}
}

func (s *IntegrationTestSuite) GetRecordExpiryQueue() {
	val := s.network.Validators[0]
	sr := s.Require()
	var authorityName = "GetRecordExpiryQueue"

	testCasesForRecordsExpiry := []struct {
		name        string
		args        []string
		expErr      bool
		noOfRecords int
		preRun      func(authorityName string, s *IntegrationTestSuite)
	}{
		{
			"invalid request",
			[]string{"invalid", fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			true,
			0,
			func(authorityName string, s *IntegrationTestSuite) {

			},
		},
		{
			"get expiry records ",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			false,
			1,
			func(authorityName string, s *IntegrationTestSuite) {
				createNameRecord(authorityName, s)
			},
		},
	}

	for _, tc := range testCasesForRecordsExpiry {
		s.Run(tc.name, func() {
			if !tc.expErr {
				tc.preRun(authorityName, s)
				time.Sleep(time.Second * 7)
			}
			cmd := cli.GetRecordExpiryQueue()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var response types.QueryGetRecordExpiryQueueResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &response)
				sr.NoError(err)
				sr.Equal(tc.noOfRecords, len(response.GetRecords()))
			}
		})
	}
}
func (s *IntegrationTestSuite) TestGetAuthorityExpiryQueue() {
	val := s.network.Validators[0]
	sr := s.Require()
	var authorityName = "TestGetAuthorityExpiryQueue"

	testCases := []struct {
		name   string
		args   []string
		expErr bool
		preRun func(authorityName string)
	}{
		{
			"invalid request without name",
			[]string{"invalid", fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			true,
			func(authorityName string) {

			},
		},
		{
			"success query",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			false,
			func(authorityName string) {
				// reserving the name
				clientCtx := val.ClientCtx
				cmd := cli.GetCmdReserveName()
				args := []string{
					authorityName,
					fmt.Sprintf("--owner=%s", accountAddress),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
					fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
					fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
					fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
				}
				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if !tc.expErr {
				tc.preRun(authorityName)
				time.Sleep(time.Second * 6)
			}
			cmd := cli.GetAuthorityExpiryQueue()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var response types.QueryGetAuthorityExpiryQueueResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &response)
				sr.NoError(err)
				sr.NotZero(len(response.GetAuthorities()))
			}
		})
	}
}

func createNameRecord(authorityName string, s *IntegrationTestSuite) {
	val := s.network.Validators[0]
	sr := s.Require()
	// reserving the name
	clientCtx := val.ClientCtx
	cmd := cli.GetCmdReserveName()
	args := []string{
		authorityName,
		fmt.Sprintf("--owner=%s", accountAddress),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
	}
	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
	sr.NoError(err)
	var d sdk.TxResponse
	err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
	sr.NoError(err)
	sr.Zero(d.Code)

	// creating the bond
	CreateBond(s)

	// Get the bond-id
	bondId := GetBondId(s)

	// adding bond-id to name authority
	args = []string{
		authorityName, bondId,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
	}
	cmd = cli.GetCmdSetAuthorityBond()

	out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
	sr.NoError(err)
	err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
	sr.NoError(err)
	sr.Zero(d.Code)

	args = []string{
		fmt.Sprintf("wrn://%s/", authorityName),
		"test_hello_cid",
		fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
	}

	cmd = cli.GetCmdSetName()

	out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
	sr.NoError(err)
	err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
	sr.NoError(err)
	sr.Zero(d.Code)
}
