cosine.Zerobased = false
sine.Zerobased = false
cosine.Frame(canvas, 5)
sine.Label(canvas, 1.5, 10)
cosine.YAxis(canvas, 1.2, -1.0, 1.0, 1.0, "%0.2f", true)
cosine.Color = color.RGBA{0, 128, 0, 255}
sine.Color = color.RGBA{128, 0, 0, 255}
cosine.Scatter(canvas, 0.5)
sine.Scatter(canvas, 0.5)
