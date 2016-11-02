package log

func (_ nullLog) Log(_ Severity, _ interface{}) {}
func (_ nullLog) Debug(_ interface{})           {}
func (_ nullLog) Info(_ interface{})            {}
func (_ nullLog) Notice(_ interface{})          {}
func (_ nullLog) Warning(_ interface{})         {}
func (_ nullLog) Error(_ interface{})           {}
func (_ nullLog) Critical(_ interface{})        {}
func (_ nullLog) Alert(_ interface{})           {}
func (_ nullLog) Emergency(_ interface{})       {}
