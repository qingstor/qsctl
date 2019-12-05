package task

//go:generate go run ../internal/cmd/generator/tasks
//go:generate mockgen -package mock -destination ../pkg/mock/storager.go github.com/Xuanwo/storage Storager,Servicer,Segmenter,Reacher
//go:generate mockgen -package mock -destination ../pkg/mock/io.go io ReadCloser
