export function getStorage(key) {
  try {
    return uni.getStorageSync(key)
  } catch (e) {
    return null
  }
}

export function setStorage(key, value) {
  uni.setStorageSync(key, value)
}

export function removeStorage(key) {
  uni.removeStorageSync(key)
}

