module github.com/liyanbing/filter

go 1.12

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20200323165209-0ec3e9974c59
	golang.org/x/mod => github.com/golang/mod v0.2.0
	golang.org/x/net => github.com/golang/net v0.0.0-20200324143707-d3edc9973b7e
	golang.org/x/sync => github.com/golang/sync v0.0.0-20200317015054-43a5402ce75a
	golang.org/x/sys => github.com/golang/sys v0.0.0-20200331124033-c3d80250170d
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/tools => github.com/golang/tools v0.0.0-20200331202046-9d5940d49312
	golang.org/x/xerrors => github.com/golang/xerrors v0.0.0-20191204190536-9bdfabe68543
)

require (
	github.com/liyanbing/calc v0.0.0-20200615034323-073f8dc291b4
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/oschwald/geoip2-golang v1.4.0
	github.com/stretchr/testify v1.6.1
)
