# Code converter
- 그래프 형식으로 그린 신경망을 파이썬 코드로 변환시키는 모듈

## 실행 예시
- 클라이언트로 부터 받은 신경망 정보 JSON 파일

```json
[
  {
    "type": "conv2d",ㅎ
    "name": "conv2d_1",
    "input": null,
    "config": {
      "shape": [
        1,
        256,
        256
      ],
      "kernel_size": [
        3,
        3
      ],
      "filter": 32,
      "stride": 1,
      "padding": "same"
    }
  },
  {
    "type": "relu",
    "name": "relu_1",
    "input": "conv2d_1",
    "config": {
    }
  },
  {
    "type": "dense",
    "name": "dense_1",
    "input": "relu_1",
    "config": {
      "units": 64
    }
  },
  {
    "type": "sigmoid",
    "name": "sigmoid_1",
    "input": "dense_1",
    "config": {
    }
  }
]
```
- 변환된 Python 코드
```python
conv2d_1 = tf.keras.layers.Conv2D(32, (3, 3), input_shape=(1, 256, 256), padding='same')(data)
relu_1 = tf.keras.layers.Activation('relu')(conv2d_1)
dense_1 = tf.keras.layers.Dense(units=64)(relu_1)
sigmoid_1 = tf.keras.layers.Activation('sigmoid)(dense_1)
```