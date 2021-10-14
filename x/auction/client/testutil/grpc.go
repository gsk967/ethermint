package testutil

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/rest"
	auctiontypes "github.com/tharsis/ethermint/x/auction/types"
)

const (
	randomAuctionID     = "randomAuctionID"
	randomBidderAddress = "randomBidderAddress"
	randomOwnerAddress  = "randomOwnerAddress"
)

func (suite *IntegrationTestSuite) TestGetAllAuctionsGrpc() {
	val := suite.network.Validators[0]
	sr := suite.Require()
	reqUrl := fmt.Sprintf("%s/vulcanize/auction/v1beta1/auctions", val.APIAddress)

	testCases := []struct {
		msg             string
		url             string
		errorMsg        string
		isErrorExpected bool
	}{
		{
			"invalid request to get all auctions",
			reqUrl + randomAuctionID,
			"",
			true,
		},
		{
			"valid request to get all auctions",
			reqUrl,
			"",
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			resp, err := rest.GetRequest(tc.url)
			if tc.isErrorExpected {
				sr.Contains(string(resp), tc.errorMsg)
			} else {
				sr.NoError(err)
				var auctions auctiontypes.AuctionsResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(resp, &auctions)
				sr.NotZero(len(auctions.Auctions.Auctions))
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestQueryParamsGrpc() {
	val := suite.network.Validators[0]
	sr := suite.Require()
	reqUrl := fmt.Sprintf("%s/vulcanize/auction/v1beta1/params", val.APIAddress)

	suite.Run("valid request to get auction params", func() {
		resp, err := rest.GetRequest(reqUrl)
		suite.Require().NoError(err)

		var params auctiontypes.QueryParamsResponse
		err = val.ClientCtx.Codec.UnmarshalJSON(resp, &params)

		sr.NoError(err)
		sr.Equal(*params.GetParams(), auctiontypes.DefaultParams())
	})
}

func (suite *IntegrationTestSuite) TestGetAuctionGrpc() {
	val := suite.network.Validators[0]
	sr := suite.Require()
	reqUrl := fmt.Sprintf("%s/vulcanize/auction/v1beta1/auctions/", val.APIAddress)

	testCases := []struct {
		msg             string
		url             string
		errorMsg        string
		isErrorExpected bool
		preRun          func() string
	}{
		{
			"invalid request to get an auction",
			reqUrl + randomAuctionID,
			"",
			true,
			func() string { return "" },
		},
		{
			"valid request to get an auction",
			reqUrl,
			"",
			false,
			func() string { return suite.defaultAuctionID },
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			auctionID := tc.preRun()
			resp, err := rest.GetRequest(tc.url + auctionID)
			if tc.isErrorExpected {
				sr.Contains(string(resp), tc.errorMsg)
			} else {
				sr.NoError(err)
				var auction auctiontypes.AuctionResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(resp, &auction)
				sr.NoError(err)
				sr.Equal(auctionID, auction.Auction.Id)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetBidsGrpc() {
	val := suite.network.Validators[0]
	sr := suite.Require()
	reqUrl := fmt.Sprintf("%s/vulcanize/auction/v1beta1/bids/", val.APIAddress)
	testCases := []struct {
		msg             string
		url             string
		errorMsg        string
		isErrorExpected bool
		preRun          func() string
	}{
		{
			"invalid request to get all bids",
			reqUrl,
			"",
			true,
			func() string { return "" },
		},
		{
			"valid request to get all bids",
			reqUrl,
			"",
			false,
			func() string { return suite.createAuctionAndBid(false, true) },
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			auctionID := tc.preRun()
			tc.url += auctionID
			resp, err := rest.GetRequest(tc.url)
			if tc.isErrorExpected {
				sr.Contains(string(resp), tc.errorMsg)
			} else {
				sr.NoError(err)
				var bids auctiontypes.BidsResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(resp, &bids)
				sr.NoError(err)
				sr.Equal(auctionID, bids.Bids[0].AuctionId)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetBidGrpc() {
	val := suite.network.Validators[0]
	sr := suite.Require()
	reqUrl := fmt.Sprintf("%s/vulcanize/auction/v1beta1/bids/", val.APIAddress)
	testCases := []struct {
		msg             string
		url             string
		errorMsg        string
		isErrorExpected bool
	}{
		{
			"invalid request to get bid",
			fmt.Sprintf("%s/%s/", reqUrl, randomAuctionID),
			"",
			true,
		},
		{
			"valid request to get bid",
			fmt.Sprintf("%s/%s/%s", reqUrl, randomAuctionID, randomBidderAddress),
			"",
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			resp, err := rest.GetRequest(tc.url)
			if tc.isErrorExpected {
				sr.Contains(string(resp), tc.errorMsg)
			} else {
				sr.NoError(err)
				var bid auctiontypes.BidResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(resp, &bid)
				sr.NoError(err)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetAuctionsByOwnerGrpc() {
	val := suite.network.Validators[0]
	sr := suite.Require()
	reqUrl := fmt.Sprintf("%s/vulcanize/auction/v1beta1/by-owner/", val.APIAddress)
	testCases := []struct {
		msg             string
		url             string
		errorMsg        string
		isErrorExpected bool
	}{
		{
			"invalid request to get auctions by owner",
			reqUrl,
			"",
			true,
		},
		{
			"valid request to get auctions by owner",
			fmt.Sprintf("%s/%s", reqUrl, randomOwnerAddress),
			"",
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			resp, err := rest.GetRequest(tc.url)
			if tc.isErrorExpected {
				sr.Contains(string(resp), tc.errorMsg)
			} else {
				sr.NoError(err)
				var auctions auctiontypes.AuctionsResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(resp, &auctions)
				sr.NoError(err)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestQueryBalanceGrpc() {
	val := suite.network.Validators[0]
	sr := suite.Require()
	reqUrl := fmt.Sprintf("%s/vulcanize/auction/v1beta1/balance", val.APIAddress)
	msg := "valid request to get the auction module balance"

	suite.createAuctionAndBid(false, true)

	suite.Run(msg, func() {
		resp, err := rest.GetRequest(reqUrl)
		sr.NoError(err)

		var response auctiontypes.BalanceResponse
		err = val.ClientCtx.Codec.UnmarshalJSON(resp, &response)

		sr.NoError(err)
		sr.NotZero(len(response.GetBalance()))
	})
}
