package service

import "github.com/sirupsen/logrus"

func LogTransformation(t Transformation, fields logrus.Fields) {
	fields["type"] = t.GetType()
	logrus.WithFields(fields).Info("New transformation")
}
