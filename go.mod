module hz.tools/fftw

go 1.15

require (
	github.com/stretchr/testify v1.7.0
	hz.tools/rf v0.0.5
	hz.tools/sdr v0.0.0-20200924134717-aca41b35cb56
)

replace hz.tools/sdr => ../sdr
