package service

import (
	"github.com/sirupsen/logrus"
)

func ValidateToken(token string) error {
	return nil
}

const (
	rootTokenExists = `Root token already generated, if you don't have it anymore
	you will have to remove it manually in order to generate it`

	rootTokenGenerated = `Root token generated: %s Please keep it somewhere safe.
  You can use it to generate new ones to allow clients to connect and use me`
)

func InitRootToken(ts TokenStore) (*Token, error) {
	rootToken, err := ts.Save(ts.RootToken())
	if err != nil && err != ErrResourceAlreadyExists {
		return nil, err
	} else {
		if err == ErrResourceAlreadyExists {
			logrus.Info(rootTokenExists)
		} else {
			logrus.Warnf(rootTokenGenerated, rootToken.Value.String())
		}
	}
	return rootToken, nil
}
