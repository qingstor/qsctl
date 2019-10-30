package task

//go:generate rm *generated*.go -f || true
//go:generate go run ../internal/cmd/generator/tasks
//go:generate mockgen -package mock -destination ../pkg/mock/storager.go github.com/Xuanwo/storage Storager,Servicer
//go:generate mockgen -package mock -destination ../pkg/mock/io.go io ReadCloser
