# Code converter
- 그래프 형식으로 그린 신경망을 파이썬 코드로 변환시키는 모듈.

#### v1.1 (2021-08-17)
- Optimizer 종류 추가
#### v1.11 (2021-08-17)
- Rescaling, Reshape 레이어 추가
- 모델 학습코드 파일 별도 생성하도록 수정

# Supports
### [Tensorflow-Keras](https://www.tensorflow.org/?hl=ko)
### Layer
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
  - Rescaling
    - scale - float
    - offset - float : default - 0.0
  - Reshape
    - target_shape - integer tuple


### Optimizer
- Adadelta
    - learning_rate - float
    - weight_decay - float
    - epsilon - float
- Adagrad
    - learning_rate - float
    - initial_accumulator_value - float
    - epsilon - float
- Adam
    - learning_rate - float
    - beta_1 - float
    - beta_2 - float
    - epsilon - float
    - amsgrad - boolean : default - false
- Adamax
    - learning_rate - float
    - beta_1 - float
    - beta_2 - float
    - epsilon - float
- Nadam
    - learning_rate - float
    - beta_1 - float
    - beta_2 - float
    - epsilon - float
- RMSprop
    - learning_rate - float
    - decay - float
    - momentum - float
    - epsilon - float
    - centered - boolean : default - false
- SGD
    - learning_rate - float
    - momentum - float
    - nesterov - boolean : default - false
- AdamW
    - weight_decay - float
    - learning_rate - float
    - beta_1 - float
    - beta_2 - float
    - epsilon - float
    - amsgrad - boolean : default - false
- SGDW
    - weight_decay - float
    - learning_rate - float
    - momentum - float
    - nesterov - boolean : default - false


## 실행 예시
- 클라이언트로 부터 받은 신경망 정보 JSON 파일

### Request body

```json
{
  "config": {
    "optimizer_name": "Adam",
    "optimizer_config": {
      "learning_rate": 0.001,
      "beta_1": 0.9,
      "beta_2": 0.999,
      "epsilon": 1e-07,
      "amsgrad": false
    },
    "loss": "sparse_categorical_crossentropy",
    "metrics": [
      "accuracy"
    ],
    "batch_size": 32,
    "epochs": 10,
    "early_stop": {
      "usage": true,
      "monitor": "loss",
      "patience": 2
    },
    "learning_rate_reduction": {
      "usage": true,
      "monitor": "val_accuracy",
      "patience": 2,
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
        "param": {
          "shape": [
            28,
            28,
            1
          ]
        }
      },
      {
        "category": "Layer",
        "type": "Conv2D",
        "name": "node_2fbbd8e5b0a5456faa2d47f7026b139f",
        "input": "node_1605430f35f94411aaf6b97eae005e19",
        "output": "node_39ce8c39bacb4fb392c2372fb81a0b7e",
        "param": {
          "filters": 16,
          "kernel_size": [
            16,
            16
          ],
          "padding": "same",
          "strides": [
            1,
            1
          ]
        }
      },
      {
        "category": "Layer",
        "type": "Dropout",
        "name": "node_2c8a6d78d0204888942f16317f2a079f",
        "input": "node_39ce8c39bacb4fb392c2372fb81a0b7e",
        "output": "node_71914b8774b64700b38dc3e8e7a62caa",
        "param": {
          "rate": 0.5
        }
      },
      {
        "category": "Layer",
        "type": "Activation",
        "name": "node_39ce8c39bacb4fb392c2372fb81a0b7e",
        "input": "node_2fbbd8e5b0a5456faa2d47f7026b139f",
        "output": "node_2c8a6d78d0204888942f16317f2a079f",
        "param": {
          "activation": "relu"
        }
      },
      {
        "category": "Layer",
        "type": "Flatten",
        "name": "node_71914b8774b64700b38dc3e8e7a62caa",
        "input": "node_2c8a6d78d0204888942f16317f2a079f",
        "output": "node_020cdce94de241ac9556bb0b0022c1f2",
        "param": {}
      },
      {
        "category": "Layer",
        "type": "Dense",
        "name": "node_020cdce94de241ac9556bb0b0022c1f2",
        "input": "node_71914b8774b64700b38dc3e8e7a62caa",
        "output": "node_96afcbc0a4ba4ed9b02b579068f166f0",
        "param": {
          "units": 10
        }
      },
      {
        "category": "Layer",
        "type": "Activation",
        "name": "node_96afcbc0a4ba4ed9b02b579068f166f0",
        "input": "node_020cdce94de241ac9556bb0b0022c1f2",
        "output": null,
        "param": {
          "activation": "softmax"
        }
      }
    ]
  }
}
```

### 변환된 Python 코드
- /make-python - model.py
```python
import tensorflow as tf

import tensorflow_addons as tfa

node_1605430f35f94411aaf6b97eae005e19 = tf.keras.layers.Input(shape=(28, 28, 1))
node_2fbbd8e5b0a5456faa2d47f7026b139f = tf.keras.layers.Conv2D(filters=16, kernel_size=(16, 16), strides=(1, 1), padding='same')(node_1605430f35f94411aaf6b97eae005e19)
node_39ce8c39bacb4fb392c2372fb81a0b7e = tf.keras.layers.Activation(activation="relu")(node_2fbbd8e5b0a5456faa2d47f7026b139f)
node_2c8a6d78d0204888942f16317f2a079f = tf.keras.layers.Dropout(rate=0.5)(node_39ce8c39bacb4fb392c2372fb81a0b7e)
node_71914b8774b64700b38dc3e8e7a62caa = tf.keras.layers.Flatten()(node_2c8a6d78d0204888942f16317f2a079f)
node_020cdce94de241ac9556bb0b0022c1f2 = tf.keras.layers.Dense(units=10)(node_71914b8774b64700b38dc3e8e7a62caa)
node_96afcbc0a4ba4ed9b02b579068f166f0 = tf.keras.layers.Activation(activation="softmax")(node_020cdce94de241ac9556bb0b0022c1f2)
model = tf.keras.Model(inputs=node_1605430f35f94411aaf6b97eae005e19, outputs=node_96afcbc0a4ba4ed9b02b579068f166f0)

model.compile(optimizer=tf.keras.optimizers.Adam(learning_rate=0.001, beta_1=0.9, beta_2=0.999, epsilon=1e-07, amsgrad=False), loss="sparse_categorical_crossentropy", metrics=["accuracy"])
```
- /fit - train.py
```python
import tensorflow as tf

import tensorflow_addons as tfa

import model


# Callback functions are below if use them.
early_stop = tf.keras.callbacks.EarlyStopping(monitor='loss', patience=2)
learning_rate_reduction = tf.keras.callbacks.ReduceLROnPlateau(monitor='val_accuracy', patience=2, verbose=1, factor=0.25, min_lr=3e-07)
remote_monitor = tf.keras.callbacks.RemoteMonitor(root='http://localohst:8080', path='/publish/epoch/end', field='data', headers=None, send_as_json=True)

model.model.fit(data, label, epochs=10, batch_size=32, validation_split=0.3, callbacks=[remote_monitor, learning_rate_reduction, early_stop])

```

> Requestbody로 넘어오는 레이어의 순서와 실제 순서가 다를 수 있기 때문에 서버에서 한 번 정렬해 준 후 코드로 변환된다.
