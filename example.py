import tensorflow as tf

input_1 = tf.keras.layers.Input(shape=(1, 26, 26))
conv2d_1 = tf.keras.layers.Conv2D(filters=32, kernel_size=(3, 3), padding='same')(input_1)
relu_1 = tf.keras.layers.Activation(activation='relu')(conv2d_1)
dense_1 = tf.keras.layers.Dense(units=64)(relu_1)
sigmoid_1 = tf.keras.layers.Activation(activation='sigmoid')(dense_1)

model = tf.keras.Model(inputs=input_1, outputs=sigmoid_1)
model.compile(optimizer='adam', loss='categorical_crossentropy', metrics=['accuracy'])
#model.fit(x=data[x], y=data[y],batch_size=32, epochs=10)