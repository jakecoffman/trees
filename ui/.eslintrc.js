module.exports = {
  root: true,
  env: {
    node: true
  },
  plugins: ['vue'],
  extends: [
    'plugin:vue/vue3-essential',
    'eslint:recommended',
    // "plugin:vue/vue3-strongly-recommended"
  ],
  parserOptions: {
    ecmaVersion: 2020
  },
  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'vue/max-attributes-per-line': 0,
    'vue/html-closing-bracket-spacing': 0
  }
}
