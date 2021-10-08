package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"net/url"
)

var (
	_ sdk.Msg = &MsgSetName{}
	_ sdk.Msg = &MsgReserveAuthority{}
	_ sdk.Msg = &MsgSetAuthorityBond{}
	_ sdk.Msg = &MsgDeleteNameAuthority{}
)

// NewMsgSetName is the constructor function for MsgSetName.
func NewMsgSetName(wrn string, cid string, signer sdk.AccAddress) MsgSetName {
	return MsgSetName{
		Wrn:    wrn,
		Cid:    cid,
		Signer: signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgSetName) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSetName) Type() string { return "set-name" }

// ValidateBasic Implements Msg.
func (msg MsgSetName) ValidateBasic() error {

	if msg.Wrn == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "WRN is required.")
	}

	if msg.Cid == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "CID is required.")
	}

	if len(msg.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer")
	}

	return nil
}

// GetSignBytes gets the sign bytes for the msg MsgSetName
func (msg MsgSetName) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgSetName) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// NewMsgReserveAuthority is the constructor function for MsgReserveName.
func NewMsgReserveAuthority(name string, signer sdk.AccAddress, owner sdk.AccAddress) MsgReserveAuthority {
	return MsgReserveAuthority{
		Name:   name,
		Owner:  owner.String(),
		Signer: signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgReserveAuthority) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgReserveAuthority) Type() string { return "reserve-authority" }

// ValidateBasic Implements Msg.
func (msg MsgReserveAuthority) ValidateBasic() error {

	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "name is required.")
	}

	if len(msg.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer")
	}

	return nil
}

// GetSignBytes gets the sign bytes for the msg MsgSetName
func (msg MsgReserveAuthority) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgReserveAuthority) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// NewMsgSetAuthorityBond is the constructor function for MsgSetAuthorityBond.
func NewMsgSetAuthorityBond(name string, bondId string, signer sdk.AccAddress) MsgSetAuthorityBond {
	return MsgSetAuthorityBond{
		Name:   name,
		Signer: signer.String(),
		BondId: bondId,
	}
}

// Route Implements Msg.
func (msg MsgSetAuthorityBond) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSetAuthorityBond) Type() string { return "authority-bond" }

// ValidateBasic Implements Msg.
func (msg MsgSetAuthorityBond) ValidateBasic() error {
	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "name is required.")
	}

	if len(msg.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer.")
	}

	if len(msg.BondId) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "bond id is required.")
	}

	return nil
}

// GetSignBytes gets the sign bytes for the msg MsgSetName
func (msg MsgSetAuthorityBond) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgSetAuthorityBond) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// NewMsgDeleteNameAuthority is the constructor function for MsgDeleteNameAuthority.
func NewMsgDeleteNameAuthority(wrn string, signer sdk.AccAddress) MsgDeleteNameAuthority {
	return MsgDeleteNameAuthority{
		Wrn:    wrn,
		Signer: signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgDeleteNameAuthority) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDeleteNameAuthority) Type() string { return "delete-name" }

// ValidateBasic Implements Msg.
func (msg MsgDeleteNameAuthority) ValidateBasic() error {
	if len(msg.Wrn) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "wrn is required.")
	}

	if len(msg.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer.")
	}

	_, err := url.Parse(msg.Wrn)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid wrn.")
	}

	return nil
}

// GetSignBytes gets the sign bytes for the msg MsgSetName
func (msg MsgDeleteNameAuthority) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgDeleteNameAuthority) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}
