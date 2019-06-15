package service

type TransformationType string

const (
	TransformationTypeRotate TransformationType = "rotate"
	TransformationTypeBlur   TransformationType = "blur"
	TransformationTypeCrop   TransformationType = "crop"
	TransformationTypeResize TransformationType = "resize"
)
