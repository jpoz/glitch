# glitch

[![GoDoc](https://godoc.org/github.com/jpoz/glitch?status.svg)](http://godoc.org/github.com/jpoz/glitch)

Glitch Images with Go (Golang)

## Functions

### NewGlitch

```go
gl, err := glitch.NewGlitch("./example.jpg")
if err != nil {
  panic(err)
}
```

> ./example.jpg

![example](https://cloud.githubusercontent.com/assets/12866/21326145/3d12fd80-c5de-11e6-8279-4a86afe26383.jpg)

---

### Seed

Glitch is based off of the `rand` package. Everything parameter is a random
value. To seed the random number generator use `Seed`.

```go
gl.Seed(10)
```

---

### Copy

Copies Input directly to Output with no manipulation. This is a good function
to start with

```go
gl.Copy()
```
![copy](https://cloud.githubusercontent.com/assets/12866/21326141/3c0c1002-c5de-11e6-943b-005697acf0c1.png)

---

### TransposeInput

```go
gl.TransposeInput()
```

![transposeinput](https://cloud.githubusercontent.com/assets/12866/21326152/3f957a4c-c5de-11e6-86a3-f268b4ccf5aa.png)

---

### VerticalTransposeInput

VerticalTransposeInput will take random vertical slices of the input, shift them
and write them to the output.

```go
gl.VerticalTransposeInput()
```

![verticaltransposeinput](https://cloud.githubusercontent.com/assets/12866/21326155/4097870a-c5de-11e6-98d4-a072545eacd9.png)

---

### CompressionGhost

```go
gl.Copy()
gl.CompressionGhost()
```

![compressionghost](https://cloud.githubusercontent.com/assets/12866/21326139/3b0eead0-c5de-11e6-8b65-c360da80c59c.png)


---

### GhostStreach

```go
gl.Copy()
gl.GhostStreach()
```

![ghoststreach](https://cloud.githubusercontent.com/assets/12866/21326148/3e4dbb40-c5de-11e6-93c9-3c65dd739919.png)


---

### HalfLifeLeft

```go
gl.Copy()
gl.HalfLifeLeft()
```


![halflifeleft](https://cloud.githubusercontent.com/assets/12866/21326319/f2b20938-c5de-11e6-838a-c3f4ad1ac571.png)

---

### HalfLifeRight

```go
gl.Copy()
gl.HalfLifeRight()
```

![halfliferight](https://cloud.githubusercontent.com/assets/12866/21326321/f3f9af9e-c5de-11e6-8d06-f65a5ad032aa.png)

---

### ChannelShiftLeft

```go
gl.Copy()
gl.ChannelShiftLeft()
```


![channelshiftleft](https://cloud.githubusercontent.com/assets/12866/21326313/ed0142a6-c5de-11e6-929f-bf46f91bbc39.png)

---

### ChannelShiftRight

```go
gl.Copy()
gl.ChannelShiftRight()
```

![channelshiftright](https://cloud.githubusercontent.com/assets/12866/21326314/ee383d32-c5de-11e6-8c3d-fff32d2a97a1.png)

---

### BlueBoost

```go
gl.Copy()
gl.BlueBoost()
```

![blueboost](https://cloud.githubusercontent.com/assets/12866/21365546/999dc5ea-c6aa-11e6-886f-f773cea07404.png)

---

### GreenBoost

```go
gl.Copy()
gl.GreenBoost()
```

![greenboost](https://cloud.githubusercontent.com/assets/12866/21365549/9b4d31f0-c6aa-11e6-839e-44d57fa6dcea.png)

---

### RedBoost

```go
gl.Copy()
gl.RedBoost()
```

![redboost](https://cloud.githubusercontent.com/assets/12866/21365554/9d56cc18-c6aa-11e6-97fa-d5d934af59b8.png)

---

### PrismBurst


```go
gl.Copy()
gl.PrismBurst()
```

![prismburst](https://cloud.githubusercontent.com/assets/12866/21365553/9d419a00-c6aa-11e6-811c-7979399445ba.png)

---

### Noise

```go
gl.Copy()
gl.Noise()
```

![noise](https://cloud.githubusercontent.com/assets/12866/21365812/dc258500-c6ab-11e6-9c99-edd1d0d7b040.png)


