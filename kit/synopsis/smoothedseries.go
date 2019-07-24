package synopsis

//SmoothedSeries compresses time series data by representing it as a sequence of averages
type SmoothedSeries struct {
	length int
	data   []float64
	caps   []int
	pos    int
	fill   int
}

//NewSmoothedSeries is a constructor for SmoothedSeries
func NewSmoothedSeries(n int, capFn func(int) int) *SmoothedSeries {
	data := make([]float64, n)
	caps := make([]int, n)
	for index := range caps {
		caps[index] = capFn(index)
	}
	return &SmoothedSeries{length: n, data: data, caps: caps}
}

//Insert inserts a sequence of values into the SmoothedSeries, in the order they are provided.
func (s *SmoothedSeries) Insert(values ...float64) {
	for _, v := range values {
		if s.pos < s.length { //Series is not at capacity yet
			pWt := s.partialWeight()
			if s.pos > 0 { //Series is past first element
				s.data[s.pos] = (1-pWt)*s.data[s.pos] + pWt*s.data[s.pos-1]
				for index := s.pos - 1; index > 0; index-- {
					wt := s.weight(index)
					s.data[index] = (1-wt)*s.data[index] + wt*s.data[index-1]
				}
				wt := s.weight(0)
				s.data[0] = (1-wt)*s.data[0] + wt*v
			} else { //Series is still filling first element
				s.data[0] = (1-pWt)*s.data[0] + pWt*v
			}
			s.updatePos()
		} else { //Series is at capacity
			for index := s.length - 1; index > 0; index-- {
				wt := s.weight(index)
				s.data[index] = (1-wt)*s.data[index] + wt*s.data[index-1]
			}
			wt := s.weight(0)
			s.data[0] = (1-wt)*s.data[0] + wt*v
		}
	}
}

func (s *SmoothedSeries) updatePos() {
	s.fill++
	if s.fill == s.caps[s.pos] {
		s.fill = 0
		s.pos++
	}
}

func (s *SmoothedSeries) weight(ind int) float64 {
	return 1 / float64(s.caps[ind])
}

func (s *SmoothedSeries) partialWeight() float64 {
	return 1 / float64(s.fill+1)
}

//Mean returns the average value
func (s *SmoothedSeries) Mean() float64 {
	mean := float64(0)
	for ind := 0; ind < s.pos; ind++ {
		mean += s.weight(ind) * s.data[ind]
	}
	if s.pos < s.length {
		mean += s.partialWeight() * s.data[s.pos]
	}
	return mean
}

//SetData sets the data slice
func (s *SmoothedSeries) SetData(data []float64) {
	if len(data) == s.length {
		s.data = data
	}
	//todo: fails silently on bad input
}

//Rescale multiplies all of the data by a constant factor
func (s *SmoothedSeries) Rescale(scale float64) {
	for ind := 0; ind < s.length; ind++ {
		s.data[ind] *= scale
	}
}

//ExponentialSmoothedSeries tracks a smoothed history of length k * (2^n - 1) using k*n values
func ExponentialSmoothedSeries(n, k int) *SmoothedSeries {
	return NewSmoothedSeries(
		n*k,
		func(r int) int { return 1 << uint(r/k) }, //2 ^ floor(r/k)
	)
}
