module github.com/qingstor/qsctl/v2

go 1.14

require (
	bou.ke/monkey v1.0.2
	github.com/AlecAivazis/survey/v2 v2.0.8
	github.com/Xuanwo/go-locale v0.3.0
	github.com/aos-dev/go-service-fs v0.0.0-20201022090612-43e1d9c08087
	github.com/aos-dev/go-service-qingstor v0.0.0-20201022030147-de5a1d3d5d50
	github.com/aos-dev/go-storage/v2 v2.0.0-20201021090247-828ece82a9ec
	github.com/c-bata/go-prompt v0.2.3
	github.com/c2h5oh/datasize v0.0.0-20171227191756-4eba002a5eae
	github.com/cosiner/argv v0.1.0
	github.com/go-openapi/strfmt v0.19.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/mattn/go-tty v0.0.3 // indirect
	github.com/pkg/term v0.0.0-20200520122047-c3ffed290a03 // indirect
	github.com/qingstor/log v0.0.0-20200804082313-615256cccabc
	github.com/qingstor/noah v0.0.0-20201106035637-d815bb5e1d15
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/vbauerster/mpb/v5 v5.3.0
	golang.org/x/crypto v0.0.0-20200214034016-1d94cc7ab1c6
	golang.org/x/text v0.3.3
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/spf13/cobra v1.0.0 => github.com/Prnyself/cobra v1.0.1-0.20200814081545-b584b1cb84aa
