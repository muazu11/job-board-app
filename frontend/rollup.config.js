import svg from 'rollup-plugin-svg'

export default {
    entry: 'src/input.js',
    dest: 'build/index.js',
    plugins: [
        svg()
    ]
}