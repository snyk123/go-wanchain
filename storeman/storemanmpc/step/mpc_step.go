package step

import (
	mpcprotocol "github.com/wanchain/go-wanchain/storeman/storemanmpc/protocol"
	mpcsyslog "github.com/wanchain/go-wanchain/storeman/syslog"
	"github.com/wanchain/go-wanchain/log"
)

type MpcMessageGenerator interface {
	initialize(*[]mpcprotocol.PeerInfo, mpcprotocol.MpcResultInterface) error
	calculateResult() error
}

type BaseMpcStep struct {
	BaseStep
	messages []MpcMessageGenerator
}

func CreateBaseMpcStep(peers *[]mpcprotocol.PeerInfo, messageNum int) *BaseMpcStep {
	return &BaseMpcStep{*CreateBaseStep(peers, -1), make([]MpcMessageGenerator, messageNum)}
}

func (mpcStep *BaseMpcStep) InitStep(result mpcprotocol.MpcResultInterface) error {
	for _, message := range mpcStep.messages {
		err := message.initialize(mpcStep.peers, result)
		if err != nil {
			log.Error("BaseMpcStep, init msg fail", "err", err.Error())
			mpcsyslog.Err("BaseMpcStep, init msg fail. err:%s", err.Error())
			return err
		}
	}

	return nil
}

func (mpcStep *BaseMpcStep) FinishStep() error {
	err := mpcStep.BaseStep.FinishStep()
	if err != nil {
		return err
	}

	for _, message := range mpcStep.messages {
		err := message.calculateResult()
		if err != nil {
			log.Error("BaseMpcStep, calculate msg result fail", "err", err.Error())
			mpcsyslog.Err("BaseMpcStep, calculate msg result fail. err:%s", err.Error())
			return err
		}
	}

	return nil
}