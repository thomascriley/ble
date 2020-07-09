package hci

import (
	"bytes"
	"encoding/binary"
	"github.com/thomascriley/ble/linux/hci/evt"
	"github.com/thomascriley/ble/linux/smp"
)

// SMP ...
type SMP interface {
	Code() int
	Marshal() []byte
	Unmarshal([]byte) error
}

func (c *Conn) sendSMP(p pdu) error {
	buf := bytes.NewBuffer(make([]byte, 0))
	if err := binary.Write(buf, binary.LittleEndian, uint16(4+len(p))); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, cidSMP); err != nil {
		return err
	}
	if err := binary.Write(buf, binary.LittleEndian, p); err != nil {
		return err
	}
	_, err := c.writePDU(buf.Bytes())
	//logger.Debug("smp", "send", fmt.Sprintf("[%X]", buf.Bytes()))
	return err
}

func (c *Conn) handleSMP(p pdu) error {
	//logger.Debug("smp", "recv", fmt.Sprintf("[%X]", p))
	code := p[0]
	//	payload := p[1:]
	/*
		var (
			resp SMP
			err  error
		)
	*/
	switch code {
	/*case smp.SecurityRequestCode:
		c.smpPairingReq = c.hci.smpCapabilites.PairingRequest()
		resp = c.smpPairingReq
	case smp.PairingRequestCode:
		c.smpPairingReq = &smp.PairingRequest{}
		if err = c.smpPairingReq.Unmarshal(payload); err == nil {
			resp, err = c.handlePairingRequest(c.smpPairingReq)
		}
	case smp.PairingResponseCode:
		c.smpPairingResp = &smp.PairingResponse{}
		if err = c.smpPairingResp.Unmarshal(p[1:]); err == nil {
			resp, err = c.handlePairingResponse(c.smpPairingResp)
		}
	case smp.PairingConfirmCode:
		confirm := &smp.PairingConfirm{}
		if err = confirm.Unmarshal(p[1:]); err == nil {
			resp, err = c.handlePairingConfirm(confirm)
		}
	case smp.PairingRandomCode:
		rand := &smp.PairingRandom{}
		if err := rand.Unmarshal(p[1:]); err != nil {
			return err
		}

		var r [16]byte
		if c.smpInitiator {
			copy(c.smpSRand[:], rand.RandomValue)
			r = c.smpSRand
		} else {
			copy(c.smpMRand[:], rand.RandomValue)
			r = c.smpMRand
		}

		confirm, err := c.generateConfirmKey(r)
		if err != nil {
			return err
		}

		// if slave calculate confirm value and check if it matches master
		check := c.smpMConfirm
		if c.smpInitiator {
			check = c.smpSConfirm
		}
		if !bytes.Equal(confirm[:], check[:]) {
			resp = &smp.PairingFailed{Reason: smp.ReasonConfirmValueFailed}
		}
		if !c.smpInitiator {
			resp = &smp.PairingRandom{RandomValue: c.smpSRand[:]}
		}

		c.ShortTermKey, err = &smp.S1(tk, c.smpSRand, c.smpMRand)
		if err != nil {
			return err
		}

		// generate IV and SKD

		// send LL_ENC_REQ

		// get LL_ENC_RSP

		//

	case smp.PairingFailedCode:
		fail := &smp.PairingFailed{}
		if err = fail.Unmarshal(p[1:]); err != nil {
			return err
		}
		return fmt.Errorf("Pairing failed: %X", fail.Reason)*/
	case smp.PairingConfirmCode, smp.PairingFailedCode, smp.PairingRandomCode, smp.PairingResponseCode, smp.PairingRequestCode, smp.SecurityRequestCode:
	case smp.EncryptionInformationCode:
	case smp.MasterIdentificationCode:
	case smp.IdentityIdentificationCode:
	case smp.IdentityAddressIdentificationCode:
	case smp.SigningInformationCode:

	case smp.PairingPublicKeyCode:
	case smp.PairingDHKeyCheckCode:
	case smp.KeypressNotificationCode:
	default:
		// If a packet is received with a reserved Code it shall be ignored. [Vol 3, Part H, 3.3]
		return nil
	}

	// FIXME: work aound to the lack of SMP implementation - always return non-supported.
	// C.5.1 Pairing Not Supported by Slave
	/*if err != nil {
		resp = &smp.PairingFailed{Reason: smp.ReasonUnspecifiedReason}
	}
	if resp == nil {*/
	resp := &smp.PairingFailed{Reason: smp.ReasonPairingNotSupported}
	//}
	return c.sendSMP(resp.Marshal())
}

func (c *Conn) handlePairingRequest(req *smp.PairingRequest) (resp SMP, err error) {
	c.smpInitiator = false
	c.smpSRand = smp.GenerateRand()
	if c.smpSConfirm, err = c.generateConfirmKey(c.smpSRand); err != nil {
		return nil, err
	}
	c.smpPairingResp = c.hci.smpCapabilites.PairingResponse(req)
	return c.smpPairingResp, nil
}

func (c *Conn) handlePairingResponse(_ *smp.PairingResponse) (resp SMP, err error) {
	c.smpInitiator = true
	c.smpMRand = smp.GenerateRand()
	if c.smpMConfirm, err = c.generateConfirmKey(c.smpMRand); err != nil {
		return nil, err
	}
	return &smp.PairingConfirm{ConfirmValue: c.smpMConfirm[:]}, nil
}

func (c *Conn) handlePairingConfirm(confirm *smp.PairingConfirm) (resp SMP, err error) {
	if !c.smpInitiator {
		copy(c.smpMConfirm[:], confirm.ConfirmValue)
		return &smp.PairingConfirm{ConfirmValue: c.smpSConfirm[:]}, nil
	}
	copy(c.smpSConfirm[:], confirm.ConfirmValue)
	return &smp.PairingRandom{RandomValue: c.smpMRand[:]}, nil
}

func (c *Conn) generateConfirmKey(rand [16]byte) ([16]byte, error) {
	method := smp.KeyGenMethodToUse(c.smpPairingReq, c.smpPairingResp, c.smpInitiator)

	tk, err := c.hci.smpCapabilites.GetTemporaryKey(method, c.smpInitiator)
	if err != nil {
		return tk, err
	}

	return smp.C1(tk, rand, c.smpPairingResp, c.smpPairingReq, c.iat(), c.ia(), c.rat(), c.ra())
}

func (c *Conn) ia() (addr [6]byte) {
	if c.smpInitiator {
		copy(addr[:], c.hci.addr)
	}
	return c.param.PeerAddress()
}

func (c *Conn) iat() byte {
	if c.smpInitiator {
		return c.hci.params.connParams.OwnAddressType
	}
	if prm, ok := c.param.(evt.LEConnectionComplete); ok {
		return prm.PeerAddressType()
	}
	return 0x00
}

func (c *Conn) ra() (addr [6]byte) {
	if c.smpInitiator {
		return c.param.PeerAddress()
	}
	copy(addr[:], c.hci.addr)
	return
}

func (c *Conn) rat() byte {
	if c.smpInitiator {
		if prm, ok := c.param.(evt.LEConnectionComplete); ok {
			return prm.PeerAddressType()
		}
		return 0x00
	}
	return c.hci.params.connParams.OwnAddressType
}
