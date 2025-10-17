<template>
  <div class="otp-input-container" @keydown="handleKeyDown">
    <input
      v-for="i in length" :key="i"
      :ref="(el) => (inputs[i - 1] = el as HTMLInputElement)"
      v-model="code[i - 1]"
      type="text"
      maxlength="1"
      class="otp-input"
      @input="handleInput(i - 1)"
      @paste="handlePaste"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, defineEmits, defineProps, defineExpose } from 'vue'

const props = defineProps({ length: { type: Number, default: 6 } })
const emit = defineEmits(['update:modelValue'])

const code = reactive(Array(props.length).fill(''))
const inputs = ref<HTMLInputElement[]>([])

const handleInput = (index: number) => {
  // Ensure only digits are entered
  if (!/^[0-9]$/.test(code[index])) {
    code[index] = ''
    return
  }

  if (code[index] && index < props.length - 1) {
    inputs.value[index + 1].focus()
  }

  emit('update:modelValue', code.join(''))
}

const handleKeyDown = (e: KeyboardEvent) => {
  const target = e.target as HTMLInputElement
  const index = inputs.value.indexOf(target)

  if (e.key === 'Backspace' && !target.value && index > 0) {
    inputs.value[index - 1].focus()
  }
}

const handlePaste = (e: ClipboardEvent) => {
  e.preventDefault()
  const pasteData = e.clipboardData?.getData('text').slice(0, props.length)
  if (pasteData) {
    for (let i = 0; i < pasteData.length; i++) {
      if (/^[0-9]$/.test(pasteData[i])) {
        code[i] = pasteData[i]
      }
    }
    emit('update:modelValue', code.join(''))
    const lastInputIndex = Math.min(pasteData.length, props.length) - 1
    if (lastInputIndex >= 0) {
        inputs.value[lastInputIndex].focus()
    }
  }
}

defineExpose({ focus: () => inputs.value[0]?.focus() })

</script>

<style scoped>
.otp-input-container {
  display: flex;
  gap: 10px;
}
.otp-input {
  width: 40px;
  height: 50px;
  text-align: center;
  font-size: 1.5rem;
  border-radius: 8px;
  border: 1px solid #dcdfe6;
  transition: border-color 0.2s;
}
.otp-input:focus {
  border-color: #409eff;
  outline: none;
}
</style>
