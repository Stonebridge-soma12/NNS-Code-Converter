import tensorflow as tf

input_1 = tf.keras.layers.Input(shape=(28, 28, 1))
conv2d_1 = tf.keras.layers.Conv2D(kernel_size=(3, 3), filters=32, strides=(1, 1), padding="same")(input_1)
activation_1 = tf.keras.layers.Activation(activation="relu")(conv2d_1)
flatten_1 = tf.keras.layers.Flatten()(activation_1)
dense_1 = tf.keras.layers.Dense(units=64)(flatten_1)
activation_2 = tf.keras.layers.Activation(activation="relu")(dense_1)
dense_2 = tf.keras.layers.Dense(units=10)(activation_2)
activation_3 = tf.keras.layers.Activation(activation="softmax")(dense_2)
model = tf.keras.Model(inputs=input_1, outputs=activation_3)

model.compile(optimizer=tf.keras.optimizers.Adam(lr=0.001000), loss="sparse_categorical_crossentropy", metrics=["accuracy"])
