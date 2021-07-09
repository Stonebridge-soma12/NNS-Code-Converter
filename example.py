import tensorflow as tf

inputs = tf.keras.Input(shape=(1, 256, 256))
conv2d_1 = tf.keras.layers.Conv2D(kernel_size=(3, 3), filters=32, strides=1, padding='same')(inputs)
relu_1 = tf.keras.layers.Activation('relu')(conv2d_1)
dense_1 = tf.keras.layers.Dense(units=64)(relu_1)
sigmoid_1 = tf.keras.layers.Activation('sigmoid')(dense_1)
model = tf.keras.Model(inputs=inputs, outputs=sigmoid_1)