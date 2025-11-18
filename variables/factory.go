package variables

import (
	"fmt"
	"strings"
	"sync"
)

func init() {
	Register(NewSimpleVariable(successVariable()))
	Register(NewSimpleVariable(randVariable()))
	Register(NewSimpleVariable(ipVariable()))
	Register(NewSimpleVariable(countryVariable()))
	Register(NewSimpleVariable(provinceVariable()))
	Register(NewSimpleVariable(cityVariable()))
	Register(NewSimpleVariable(timeStampVariable()))
	Register(NewSimpleVariable(tsSimpleVariable()))
	Register(NewSimpleVariable(secondVariable()))
	Register(NewSimpleVariable(minuteVariable()))
	Register(NewSimpleVariable(hourVariable()))
	Register(NewSimpleVariable(dayVariable()))
	Register(NewSimpleVariable(monthVariable()))
	Register(NewSimpleVariable(yearVariable()))
	Register(NewSimpleVariable(wdayVariable()))
	Register(NewSimpleVariable(dateVariable()))
	Register(NewSimpleVariable(timeVariable()))
	Register(NewSimpleVariable(uaVariable()))
	Register(NewSimpleVariable(refererVariable()))
	Register(NewSimpleVariable(isLoginVariable()))
	Register(NewSimpleVariable(versionVariable()))
	Register(NewSimpleVariable(platformVariable()))
	Register(NewSimpleVariable(channelVariable()))
	Register(NewSimpleVariable(uidVariable()))
	Register(NewSimpleVariable(deviceVariable()))
	Register(NewSimpleVariable(userTagVariable()))
	Register(newDataBuilder())
	Register(newCalcBuilder())
	Register(newFreqBuilder())
	Register(newCtxBuilder())
}

var defaultFactory = &factory{
	builder: make(map[string]Builder),
}

type factory struct {
	builder map[string]Builder
	sync.Mutex
}

func (s *factory) Get(name string) (Variable, bool) {
	if builder, ok := s.builder[name]; ok {
		return builder.Build(name), true
	}

	segments := strings.Split(name, ".")
	if len(segments) <= 0 {
		return nil, false
	}

	if builder, ok := s.builder[segments[0]+"."]; ok {
		return builder.Build(name), true
	}
	return nil, false
}

func (s *factory) Register(builder Builder) {
	s.Lock()
	defer s.Unlock()
	name := builder.Name()
	if _, ok := s.builder[name]; ok {
		panic(fmt.Sprintf("%v variable already exists", name))
	}
	s.builder[name] = builder
}

func Register(builder Builder) {
	defaultFactory.Register(builder)
}

func Get(name string) (Variable, bool) {
	return defaultFactory.Get(name)
}
