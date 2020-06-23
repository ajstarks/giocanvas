sr, err := os.Open("sin.d")
if err != nil {
	return err
}
cr, err := os.Open("cos.d")
if err != nil {
	return err
}
sine, err := chart.DataRead(sr)
if err != nil {
	return err
}
cosine, err := chart.DataRead(cr)
if err != nil {
	return err
}
