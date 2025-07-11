# Jump 'n' Shoot

> Um jogo estilo *run-and-gun* minimalista feito para o [WASM-4](https://wasm4.org/), inspirado em clássicos como **Contra** e o jogo do dinossauro do Chrome.

---

## 🎮 Sobre o Jogo

**Jump 'n' Shoot** é um jogo de plataforma com foco em ação e reflexos rápidos. O jogador deve desviar de obstáculos, eliminar inimigos e sobreviver o máximo possível. A pontuação aumenta ao destruir inimigos ou projéteis inimigos.

> ⚠️ Este é um projeto em desenvolvimento — contribuições, sugestões e correções são bem-vindas!

---

## 🕹️ Controles

| Ação    | Teclas WASM-4     |
|---------|-------------------|
| Pular   | Botão 1 (Z / C)   |
| Atirar  | Botão 2 (X / V)   |

💡 Também é possível usar o **mouse**:
- **Clique esquerdo**: Pular  
- **Clique direito**: Atirar

---

## 🚀 Como Jogar Localmente

### 1. Pré-requisitos

Certifique-se de ter os seguintes itens instalados:

- [WASM-4 SDK](https://wasm4.org/docs/getting-started/)
- [Go](https://golang.org/)
- [TinyGo](https://tinygo.org/)

> ⚙️ O jogo foi desenvolvido com **TinyGo**, que é necessário para compilar o código para WebAssembly.

### 2. Clonar o repositório

git clone https://github.com/seu-usuario/jump-n-shoot.git
cd jump-n-shoot


### 2. Compilar e Executar

```bash
make
make run
```
