# Code converter
- 그래프 형식으로 그린 신경망을 파이썬 코드로 변환시키는 모듈.

# Supports
### [Tensorflow-Keras](https://www.tensorflow.org/?hl=ko)
  - Dense
    - units - integer
  - Conv2D
    - filters - integer
    - kernel_size - tuple (n, n)
    - strides - tuple (n, n)
    - padding - string (drop다운으로 선택하면 될듯)
  - AveragePooling2D
    - pool_size - tuple (n, n)
    - strides - tuple (n, n)
    - padding - string
  - MaxPool2D
    - pool_size - tuple (n, n)
    - strides - tuple (n, n)
    - padding -string
  - Activation
    - activation - string
  - Input
    - shape - tuple
  - Dropout
    - rate - float 0 ~ 1
  - BatchNormalization
    - axis - integer
    - momentum - float
    - epsilon - float
  - Flatten


## 실행 예시
- 클라이언트로 부터 받은 신경망 정보 JSON 파일

### Request body
```json
{
  "config": {
    "optimizer": "adam",
    "learning_rate": 0.001,
    "loss": "sparse_categorical_crossentropy",
    "metrics": ["accuracy"],
    "batch_size": 32,
    "epochs": 10,
    "early_stop": {
      "usage": true,
      "monitor": "loss",
      "patience": 2,
    },
    "learning_rate_reduction": {
      "usage": true,
      "monitor": "val_accuracy",
      "patience":2,
      "factor": 0.25,
      "min_lr": 0.0000003
    }
  },
  "content": {
    "output": "node_96afcbc0a4ba4ed9b02b579068f166f0",
    "input": "node_1605430f35f94411aaf6b97eae005e19",
    "layers": [
      {
        "category": "Layer",
        "type": "Input",
        "name": "node_1605430f35f94411aaf6b97eae005e19",
        "input": null,
        "output": "node_2fbbd8e5b0a5456faa2d47f7026b139f",
        "config": {
          "shape": ""
        }
      },
      {
        "category": "Layer",
        "type": "Conv2D",
        "name": "node_2fbbd8e5b0a5456faa2d47f7026b139f",
        "input": "node_1605430f35f94411aaf6b97eae005e19",
        "output": "node_39ce8c39bacb4fb392c2372fb81a0b7e",
        "config": {
          "filters": "",
          "kernel_size": "",
          "padding": "",
          "strides": ""
        }
      },
      {
        "category": "Layer",
        "type": "Dropout",
        "name": "node_2c8a6d78d0204888942f16317f2a079f",
        "input": "node_39ce8c39bacb4fb392c2372fb81a0b7e",
        "output": "node_71914b8774b64700b38dc3e8e7a62caa",
        "config": {
          "rate": ""
        }
      },
      {
        "category": "Layer",
        "type": "Activation",
        "name": "node_39ce8c39bacb4fb392c2372fb81a0b7e",
        "input": "node_2fbbd8e5b0a5456faa2d47f7026b139f",
        "output": "node_2c8a6d78d0204888942f16317f2a079f",
        "config": {
          "activation": "relu"
        }
      },
      {
        "category": "Layer",
        "type": "Flatten",
        "name": "node_71914b8774b64700b38dc3e8e7a62caa",
        "input": "node_2c8a6d78d0204888942f16317f2a079f",
        "output": "node_020cdce94de241ac9556bb0b0022c1f2",
        "config": {}
      },
      {
        "category": "Layer",
        "type": "Dense",
        "name": "node_020cdce94de241ac9556bb0b0022c1f2",
        "input": "node_71914b8774b64700b38dc3e8e7a62caa",
        "output": "node_96afcbc0a4ba4ed9b02b579068f166f0",
        "config": {
          "units": ""
        }
      },
      {
        "category": "Layer",
        "type": "Activation",
        "name": "node_96afcbc0a4ba4ed9b02b579068f166f0",
        "input": "node_020cdce94de241ac9556bb0b0022c1f2",
        "output": null,
        "config": {
          "activation": "softmax"
        }
      }
    ]
  }
}

```

- 변환된 Python 코드
```python
import tensorflow as tf

node_1605430f35f94411aaf6b97eae005e19 = tf.keras.layers.Input(shape="")
node_2fbbd8e5b0a5456faa2d47f7026b139f = tf.keras.layers.Conv2D(filters="", kernel_size="", padding="", strides="")(node_1605430f35f94411aaf6b97eae005e19)
node_39ce8c39bacb4fb392c2372fb81a0b7e = tf.keras.layers.Activation(activation="relu")(node_2fbbd8e5b0a5456faa2d47f7026b139f)
node_2c8a6d78d0204888942f16317f2a079f = tf.keras.layers.Dropout(rate="")(node_39ce8c39bacb4fb392c2372fb81a0b7e)
node_71914b8774b64700b38dc3e8e7a62caa = tf.keras.layers.Flatten()(node_2c8a6d78d0204888942f16317f2a079f)
node_020cdce94de241ac9556bb0b0022c1f2 = tf.keras.layers.Dense(units="")(node_71914b8774b64700b38dc3e8e7a62caa)
node_96afcbc0a4ba4ed9b02b579068f166f0 = tf.keras.layers.Activation(activation="softmax")(node_020cdce94de241ac9556bb0b0022c1f2)
model = tf.keras.Model(inputs=node_1605430f35f94411aaf6b97eae005e19, outputs=node_96afcbc0a4ba4ed9b02b579068f166f0)

model.compile(optimizer=tf.keras.optimizers.adam(lr=0.001000), loss="sparse_categorical_crossentropy", metrics=["accuracy"])

```

> Requestbody로 넘어오는 레이어의 순서와 실제 순서가 다를 수 있기 때문에 서버에서 한 번 정렬해 준 후 코드로 변환된다.
