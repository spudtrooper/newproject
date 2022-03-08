package gen

//go:generate genopts --outfile=options.go "outdir:string"

type Option func(*optionImpl)

type Options interface {
	Outdir() string
}

func Outdir(outdir string) Option {
	return func(opts *optionImpl) {
		opts.outdir = outdir
	}
}

type optionImpl struct {
	outdir string
}

func (o *optionImpl) Outdir() string { return o.outdir }

func makeOptionImpl(opts ...Option) *optionImpl {
	res := &optionImpl{}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func MakeOptions(opts ...Option) Options {
	return makeOptionImpl(opts...)
}
