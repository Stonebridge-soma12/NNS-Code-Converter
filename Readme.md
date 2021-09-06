# Code converter
- 그래프 형식으로 그린 신경망을 파이썬 코드로 변환시키는 모듈.

#### v1.1 (2021-08-17)
- Optimizer 종류 추가
#### v1.11 (2021-08-17)
- Rescaling, Reshape 레이어 추가
- 모델 학습코드 파일 별도 생성하도록 수정

#### v1.20 (2021-08-23)
- Train을 Python 서버에서 실행.
    - CodeConverter으로 학습 요청이 들어오면 파이썬코드로 만들어진 모델을 save
    - save된 모델을 압축한 후 Python 서버로 Model config을 Body에 실어 Post요청
    - [Python server](https://github.com/Stonebridge-soma12/GPUServer) 는 Codeconverter 서버에 압축된 모델을 받아온다.
    - Python server에서 압축된 모델을 압축 해제 후 load_model
    - body에 딸려온 Config에서 데이터셋 정보를 갖고 데이터 가공 후 학습.
- 데이터셋 정보를 Config에 탑재 (임시)

#### v1.21 (2021-08-25)
- Math 모듈 추가

#### v1.30 (2021-09-02)
- Python server 사이에 Message queue 추가
    - 각 Trainer (Worker) 마다 하나의 학습 요청만 처리하도록 설정
    - flask 프레임워크가 필요없어졌기 때문에 삭제

#### v1.31 (2021-08-05)
- 한 노드에 여러 입출력 처리 추가
    - 레이어 변수명 먼저 모두 선언 후 Input 레이어부터 BFS로 순회하며 각 레이어 연걸.

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

### Math
- Abs
- Ceil
- Floor
- Round
- Sqrt


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
    "loss": "binary_crossentropy",
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
  "dataset": {
    "train_uri": "https://dataset",
    "valid_uri": "",
    "shuffle": false,
    "label": "blue_win",
    "normalization": {
      "usage": true,
      "method": "MinMax"
    }
  },
  "content": {
    "output": "Activation_2",
    "input": "Input_1",
    "layers": [
      {
        "category": "Layer",
        "type": "Input",
        "name": "Input_1",
        "input": null,
        "output": [
          "Dense_1"
        ],
        "param": {
          "shape": [
            1,
            58
          ]
        }
      },
      {
        "category": "Layer",
        "type": "Dense",
        "name": "Dense_1",
        "input": [
          "Input_1"
        ],
        "output": [
          "Activation_1"
        ],
        "param": {
          "units": 256
        }
      },
      {
        "category": "Layer",
        "type": "Activation",
        "name": "Activation_1",
        "input": [
          "Dense_1"
        ],
        "output": [
          "Dense_2"
        ],
        "param": {
          "activation": "relu"
        }
      },
      {
        "category": "Layer",
        "type": "Dense",
        "name": "Dense_2",
        "input": [
          "Activation_1"
        ],
        "output": [
          "Activation_2"
        ],
        "param": {
          "units": 1
        }
      },
      {
        "category": "Layer",
        "type": "Activation",
        "name": "Activation_2",
        "input": [
          "Dense_2"
        ],
        "output": null,
        "param": {
          "activation": "sigmoid"
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

Input_1 = tf.keras.layers.Input(shape=(1, 58))
Dense_1 = tf.keras.layers.Dense(units=256)
Activation_1 = tf.keras.layers.Activation(activation="relu")
Dense_2 = tf.keras.layers.Dense(units=1)
Activation_2 = tf.keras.layers.Activation(activation="sigmoid")
Dense_1 = Dense_1(Input_1)
Activation_1 = Activation_1(Dense_1)
Dense_2 = Dense_2(Activation_1)
Activation_2 = Activation_2(Dense_2)
model = tf.keras.Model(inputs=Input_1, outputs=Activation_2)

model.compile(optimizer=tf.keras.optimizers.Adam(learning_rate=0.001, beta_1=0.9, beta_2=0.999, epsilon=1e-07, amsgrad=False), loss="binary_crossentropy", metrics=["accuracy"])
```

> Requestbody로 넘어오는 레이어의 순서와 실제 순서가 다를 수 있기 때문에 서버에서 한 번 정렬해 준 후 코드로 변환된다.
