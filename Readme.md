# Code converter
- 그래프 형식으로 그린 신경망을 파이썬 코드로 변환시키는 모듈

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
    - strides - tuple (n, n) (None 가능)
    - padding - string
  - MaxPool2D
    - pool_size - tuple (n, n)
    - strides - tuple (n, n) (None 가능)
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

```json
// Project config
{
  "optimizer": "adam",
  "learning_rate": 0.001,
  "loss": "sparse_categorical_crossentropy",
  "metrics": ["accuracy"],
  "batch_size": 32,
  "epochs": 10,
}

// Content
{
  "output": "activation_3",
  "layers": [
    {
      "type": "input",
      "name": "input_1",
      "input": null,
      "config": {
        "input_shape": [28, 28, 1],
      }
    },
    {
      "type": "conv2d",
      "name": "conv2d_1",
      "input": "input_1",
      "config": {
        "kernel_size": [3, 3],
        "filters": 32,
        "stride": 1,
        "padding": "same"
      }
    },
    {
      "type": "activation",
      "name": "activation_1",
      "input": "conv2d_1",
      "config": {
        "activation": "relu",
      }
    },
    {
      "type": "flatten",
      "name": "flatten_1",
      "input": "activation_1",
      "config": {
      }
    },
    {
      "type": "dense",
      "name": "dense_1",
      "input": "flatten_1",
      "config": {
        "units": 64,
      }
    },
    {
      "type": "activation",
      "name": "activation_2",
      "input": "dense_1",
      "config": {
        "activation": "relu",
      }
    },
    {
      "type": "dense",
      "name": "dense_2",
      "input": "activation_2",
      "config": {
        "units": 10,
      }
    },
    {
      "type": "activation",
      "name": "activation_3",
      "input": "dense_2",
      "config": {
        "activation": "softmax",
      }
    },
  ],
}
```
- 변환된 Python 코드
```python
import tensorflow as tf

input_1 = tf.keras.layers.Input(shape=(28, 28, 1))
conv2d_1 = tf.keras.layers.Conv2D(filters=32, kernel_size=(3, 3), padding='same', strides=1)(input_1)
activation_1 = tf.keras.layers.Activation(activation='relu')(conv2d_1)
flatten_1 = tf.keras.layers.Flatten()(activation_1)
dense_1 = tf.keras.layers.Dense(units=64)(flatten_1)
activation_2 = tf.keras.layers.Activation(activation='relu')(dense_1)
dense_2 = tf.keras.layers.Dense(units=10)(activation_2)
activation_3 = tf.keras.layers.Activation(activation='softmax')(dense_2)

model = tf.keras.Model(inputs=input_1, outputs=activation_3)
model.compile(optimizer='adam', loss='sparse_categorical_crossentropy', metrics=['accuracy'])
model.fit(x=train, y=label,batch_size=32, epochs=10)

```