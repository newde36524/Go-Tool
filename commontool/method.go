package commontool

type InteractionLocker struct {
	left  chan int
	right chan int
	check Check
	sign  int
}

func NewInteractionLocker(check Check) *InteractionLocker {
	return &InteractionLocker{
		check: check,
		left:  make(chan int, 1),
		right: make(chan int, 1),
		sign:  0,
	}
}

type Check func(a, b int) bool

//Interaction 当前方法调用会阻塞，直到Check方法返回true
func (lock *InteractionLocker) Interaction(o int) bool {
	lock.sign = lock.sign % 2
	if lock.sign == 0 {
		select {
		case lock.left <- o:
		default:
		}
	}
	if lock.sign == 1 {
		select {
		case lock.right <- o:
		default:
		}
	}
	lock.sign++
	if lock.sign == 2 {
		var a int
		var b int
		select {
		case a = <-lock.left:
		default:
		}
		select {
		case b = <-lock.right:
		default:
		}
		lock.check(a, b)
		lock.sign = 0
	}
	return false
}
func (lock *InteractionLocker) do() {
	lock.sign++
	if lock.sign == 2 {
		var a int
		var b int
		select {
		case a = <-lock.left:
		default:
		}
		select {
		case b = <-lock.right:
		default:
		}
		lock.check(a, b)
	}
}
func (lock *InteractionLocker) Left(o int) {
	lock.left <- o
	lock.sign = lock.sign % 2
	lock.do()
}
func (lock *InteractionLocker) Right(o int) {
	lock.right <- o
	lock.sign = lock.sign % 2
	lock.do()
}
