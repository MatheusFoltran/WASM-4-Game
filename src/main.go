package main

import "unsafe"

// WASM-4 Memory-mapped I/O
var (
    DRAW_COLORS = (*uint16)(unsafe.Pointer(uintptr(0x14)))
    GAMEPAD1    = (*uint8)(unsafe.Pointer(uintptr(0x16)))
    PALETTE     = (*[4]uint32)(unsafe.Pointer(uintptr(0x04)))
    MOUSE_BUTTONS = (*uint8)(unsafe.Pointer(uintptr(0x1e)))
)

// Constantes otimizadas
const (
    SCREEN_WIDTH  = 160
    SCREEN_HEIGHT = 160
    GROUND_Y      = 130
    
    // Botões
    BUTTON_1 = 1
    BUTTON_2 = 2

    // Botões do mouse
    MOUSE_LEFT  = 1
    MOUSE_RIGHT = 2
    
    // Estados do jogo
    STATE_MENU = 0
    STATE_PLAYING = 1
    STATE_GAME_OVER = 2
    
    // Jogador
    PLAYER_WIDTH = 8
    PLAYER_HEIGHT = 12
    
    // Munição
    MAX_BULLETS = 4
    MAX_AMMO = 8
    RELOAD_TIME = 120
    BULLET_WIDTH = 4
    BULLET_HEIGHT = 2

    // Tiro inimigo
    MAX_ENEMY_BULLETS = 6
    ENEMY_BULLET_WIDTH = 3
    ENEMY_BULLET_HEIGHT = 3
    ENEMY_SHOOT_RATE = 30
    
    // Inimigos
    MAX_ENEMIES = 3
    ENEMY_GROUND = 0
    ENEMY_FLYING = 1
    
    // Obstaculos
    MAX_OBSTACLES = 3
    OBSTACLE_ROCK = 0
    OBSTACLE_SPIKE = 1

    // Constantes para mira
    AIM_HORIZONTAL = 0
    AIM_VERTICAL   = 1
    
    // Botões de direção
    BUTTON_LEFT  = 16
    BUTTON_RIGHT = 32
    BUTTON_UP    = 64
    BUTTON_DOWN  = 128
    
    // Particulas
    MAX_PARTICLES = 4

    // Procedural
    PATTERN_EASY = 0
    PATTERN_JUMP = 1
    PATTERN_SHOOT = 2
)

// Sprites
var (
    // Jogador (8x12) - dois frames: parado e correndo
    playerSprite = [2][12]uint8{
        // Frame 0 - parado
        {
            0b00011000, // cabeça
            0b00111100, // cabeça
            0b00011000, // pescoço
            0b00111100, // ombros
            0b01111110, // torso
            0b00111100, // cintura
            0b00111100, // quadril
            0b00111100, // pernas
            0b01111110, // pernas
            0b01111110, // pernas
            0b01100110, // pés
            0b11100111, // pés
        },
        // Frame 1 - correndo
        {
            0b00011000, // cabeça
            0b00111100, // cabeça
            0b00011000, // pescoço
            0b00111100, // ombros
            0b01111110, // torso
            0b00111100, // cintura
            0b00111100, // quadril
            0b00111100, // pernas
            0b01111110, // pernas
            0b01111110, // pernas
            0b11001100, // pés correndo
            0b11001100, // pés correndo
        },
    }

    // Sprite da bala para UI (4x2)
    bulletSprite = [2]uint8{
        0b11110000, // corpo da bala
        0b11110000, // corpo da bala
    }
    
    // Sprite da arma (6x3) - pistola
    weaponSprite = [3]uint8{
        0b11111100, // cano
        0b11000000, // corpo
        0b11000000, // cabo
    }

    // Sprite da arma vertical (3x6) - pistola apontando para cima
    weaponVerticalSprite = [6]uint8{
        0b11000000, // cabo
        0b11000000, // cabo
        0b11000000, // corpo
        0b11000000, // corpo
        0b11100000, // cano
        0b11100000, // cano
    }
    
    // Sprite do inimigo terrestre (8x12)
    groundEnemySprite = [12]uint8{
        0b00111100, // antenas
        0b01111110, // cabeça
        0b11100111, // olhos
        0b01111110, // cabeça
        0b11111111, // ombros
        0b01111110, // torso
        0b11111111, // torso
        0b01111110, // cintura
        0b11111111, // pernas
        0b01111110, // pernas
        0b11111111, // pés
        0b10100101, // rodas
    }
    
    // Sprite do inimigo aéreo (8x8)
    flyingEnemySprite = [8]uint8{
        0b10100101, // hélices
        0b01111110, // corpo
        0b11111111, // corpo
        0b11100111, // sensores
        0b11111111, // corpo
        0b01111110, // corpo
        0b00111100, // base
        0b10100101, // hélices
    }
    
    // Obstáculo - Pedra (8x8)
    rockSprite = [8]uint8{
        0b00111100,
        0b01111110,
        0b11111111,
        0b11111111,
        0b11111111,
        0b11111111,
        0b01111110,
        0b00111100,
    }
    
    // Obstáculo - Espeto (6x8)
    spikeSprite = [8]uint8{
        0b00100000,
        0b01110000,
        0b11111000,
        0b11111000,
        0b11111000,
        0b11111000,
        0b11111000,
        0b11111000,
    }
)

// Variável global para controlar direção da mira
var aimDirection int8 = AIM_HORIZONTAL

// Cooldown
var gameOverTimer uint8
var previousGamepadState uint8

// Variáveis globais para controle de entrada
var (
    prevMouseButtons uint8 = 0
    prevGamepad uint8 = 0
)

// Sistema de velocidade progressiva
var (
    currentPlayerSpeed int32 = 1
    currentEnemySpeed int32 = 1
    currentBulletSpeed int32 = 3
    currentJumpPower int8 = -11
)

// Sistema de geração procedural
var (
    rngSeed uint32 = 12345 // Seed inicial
    spawnTimer int32 = 0
    nextSpawnDelay int32 = 60
    lastSpawnType int8 = -1 // Controla tipos consecutivos
    lastSpawnX int32 = 0    // Controla proximidade
    currentPattern int8 = PATTERN_EASY
    patternTimer int32 = 0
    gameStartRealTime int32 = 0
    frameCounter int32 = 0
)

// Game state
var (
    gameState int8 = STATE_MENU
    gameFrame int32 = 0
    cameraX   int32 = 0
    score     int32 = 0
    highScore int32 = 0
)

// Sistema de munição
var (
    ammo int32 = MAX_AMMO
    reloadTimer int32 = 0
    isReloading bool = false
)

// Jogador
var player struct {
    x, y      int32
    velY      int8
    flags     uint8 // bit 0: onGround, bit 1: alive
    animFrame int8
}

// Tiros
var bullets [MAX_BULLETS]struct {
    x, y        int32
    velX, velY  int8
    active      bool
}

// Inimigos
var enemies [MAX_ENEMIES]struct {
    x, y        int32
    velX, velY  int8
    enemyType   int8
    active      bool
    animFrame   int8
}

// Tiros inimigos
var enemyBullets [MAX_ENEMY_BULLETS]struct {
    x, y   int32
    velX, velY int8
    active bool
}

// Obstáculos
var obstacles [MAX_OBSTACLES]struct {
    x, y        int32
    width, height int8
    obstacleType int8
    active      bool
}

// Partículas
var particles [MAX_PARTICLES]struct {
    x, y       int32
    velX, velY int8
    life       int8
    active     bool
}

//go:export start
func start() {
    PALETTE[0] = 0x1a1c2c // Azul escuro (fundo)
    PALETTE[1] = 0x1d2b53 // Verde petróleo (chão)
    PALETTE[2] = 0xab1c2f // Vermelho (inimigos/player)
    PALETTE[3] = 0xf4a261 // Laranja queimado (obstáculos)
    
    initGame()
}

//go:export update
func update() {
    frameCounter++
    gameFrame++
    
    switch gameState {
    case STATE_MENU:
        updateMenu()
    case STATE_PLAYING:
        updateGame()
    case STATE_GAME_OVER:
        updateGameOver()
    }
    
    draw()
}

// Inicialização
func initGame() {
    initPlayer()
    clearArrays()
}

// Inicia Jogador
func initPlayer() {
    player.x = 20
    player.y = GROUND_Y - PLAYER_HEIGHT
    player.velY = 0
    player.flags = 0x03 // onGround=1, alive=1
    player.animFrame = 0
}

func clearArrays() {
    for i := 0; i < MAX_BULLETS; i++ {
        bullets[i].active = false
        bullets[i].velY = 0
    }
    for i := 0; i < MAX_ENEMY_BULLETS; i++ {
        enemyBullets[i].active = false
    }
    for i := 0; i < MAX_ENEMIES; i++ {
        enemies[i].active = false
    }
    for i := 0; i < MAX_OBSTACLES; i++ {
        obstacles[i].active = false
    }
    for i := 0; i < MAX_PARTICLES; i++ {
        particles[i].active = false
    }
}

// Velocidades para níveis de dificuldade
func updateSpeeds() {
    if score >= 500 {
        // Nível extremo
        currentPlayerSpeed = 3
        currentEnemySpeed = 2
        currentBulletSpeed = 5
        currentJumpPower = -9
    } else if score >= 300 {
        // Nível difícil
        currentPlayerSpeed = 2
        currentEnemySpeed = 2
        currentBulletSpeed = 5
        currentJumpPower = -10
    } else if score >= 100 {
        // Nível médio
        currentPlayerSpeed = 2
        currentEnemySpeed = 1
        currentBulletSpeed = 3
        currentJumpPower = -10
    }
    // else: mantém velocidades iniciais (pulo mais forte para compensar velocidade baixa)
}

// Gerador de números pseudo-aleatórios simples
func rng() uint32 {
    rngSeed = (rngSeed*1664525 + 1013904223) & 0x7FFFFFFF
    // Melhora a distribuição com XOR shift
    rngSeed = rngSeed ^ (rngSeed >> 16)
    rngSeed = rngSeed ^ (rngSeed << 3)
    return rngSeed
}

// Retorna um número entre 0 e max-1
func randInt(max int32) int32 {
    // Chama rng() duas vezes para melhor distribuição
    val1 := rng()
    val2 := rng()
    combined := (val1 + val2) / 2
    return int32(combined % uint32(max))
}

// Função para embaralhar ainda mais a seed durante o jogo
func shuffleSeed() {
    rngSeed = rngSeed ^ uint32(gameFrame)
    rngSeed = (rngSeed << 7) ^ (rngSeed >> 25)
    rngSeed = rngSeed ^ uint32(score*3 + int32(player.x/10))
}

// Gera aleatoriedade no design do jogo
func proceduralSpawn() {

    // Embaralha a seed periodicamente para mais variação
    if spawnTimer % 60 == 0 {
        shuffleSeed()
    }

    spawnTimer++
    patternTimer++
    
    // Muda padrão com tempo mais longo (15-25 segundos)
    patternChangeTime := 900 + randInt(600) // 15-25 segundos
    if patternTimer > int32(patternChangeTime) {
        patternTimer = 0
        // Evita repetir o mesmo padrão
        newPattern := int8(randInt(4))
        if newPattern == currentPattern {
            newPattern = int8((newPattern + 1) % 4)
        }
        currentPattern = newPattern
    }
    
    // Re-embaralha baseado no tempo real transcorrido
    if spawnTimer % 180 == 0 { // A cada 3 segundos
        timeSeed := frameCounter + gameFrame*7 + score*11
        rngSeed = rngSeed ^ uint32(timeSeed)
        rng() // Consome um número para embaralhar mais
    }

    if spawnTimer >= nextSpawnDelay {
        spawnTimer = 0
        
        // Posição de spawn mais variada
        baseSpawnX := cameraX + SCREEN_WIDTH + 30 + randInt(80)
        
        // Evita spawns muito próximos com distância variável
        minDistance := 50 + randInt(40)
        if baseSpawnX - lastSpawnX < int32(minDistance) {
            baseSpawnX = lastSpawnX + int32(minDistance) + randInt(30)
        }
        
        // Spawn baseado no padrão atual
        switch currentPattern {
        case PATTERN_EASY:
            spawnEasyPattern(baseSpawnX)
        case PATTERN_JUMP:
            spawnJumpPattern(baseSpawnX)
        case PATTERN_SHOOT:
            spawnShootPattern(baseSpawnX)
        }
        
        lastSpawnX = baseSpawnX
        
        // Delay simples e fixo
        nextSpawnDelay = 60 + randInt(30) // Entre 60-90 frames (1-1.5 segundos)
    }
}

func spawnEasyPattern(x int32) {
    // Divide claramente: 50% inimigos, 40% obstáculos, 10% nada
    choice := randInt(100)
    if choice < 25 {
        // Só inimigo terrestre
        spawnEnemy(x, GROUND_Y-12, ENEMY_GROUND)
    } else if choice < 50 {
        // Só inimigo voador
        spawnEnemy(x, 80 + randInt(20), ENEMY_FLYING)
    } else if choice < 90 {
        // Só obstáculos
        if randInt(2) == 0 {
            spawnObstacle(x, GROUND_Y-8, 8, 8, OBSTACLE_ROCK)
        } else {
            spawnObstacle(x, GROUND_Y-8, 6, 8, OBSTACLE_SPIKE)
        }
    }
    // 10% sem nada
}

func spawnJumpPattern(x int32) {
    choice := randInt(100)
    if choice < 40 {
        spawnObstacle(x, GROUND_Y-8, 8, 8, OBSTACLE_ROCK)
    } else if choice < 60 {
        spawnObstacle(x, GROUND_Y-8, 6, 8, OBSTACLE_SPIKE)
    } else if choice < 80 {
        spawnObstacle(x, GROUND_Y-8, 8, 8, OBSTACLE_ROCK)
    } else if choice < 95 {
        obstacleType := OBSTACLE_ROCK
        if randInt(2) == 0 {
            obstacleType = OBSTACLE_SPIKE
        }
        spawnObstacle(x, GROUND_Y-8, 8, 8, int8(obstacleType))
    }
    // 5% sem nada
}

func spawnShootPattern(x int32) {
    choice := randInt(100)
    if choice < 50 {
        // Um inimigo voador apenas
        spawnEnemy(x, 70 + randInt(30), ENEMY_FLYING)
    } else if choice < 80 {
        // Dois inimigos voadores em alturas diferentes
        spawnEnemy(x, 60 + randInt(20), ENEMY_FLYING)
        spawnEnemy(x + 150 + randInt(80), 90 + randInt(20), ENEMY_FLYING)
    } else {
        // Dois inimigos terrestres (mais fácil que terrestre + voador)
        spawnEnemy(x, GROUND_Y-12, ENEMY_GROUND)
        spawnEnemy(x + 120 + randInt(70), GROUND_Y-12, ENEMY_GROUND)
    }
}

func updateMenu() {
    gamepad := *GAMEPAD1
    
    if gamepad&(BUTTON_1|BUTTON_2) != 0 {
        gameState = STATE_PLAYING
        startGame()
    }
}

func updateGame() {
    handleInput()
    updateAmmo()
    updateSpeeds()
    updatePlayer()
    updateBullets()
    updateEnemyBullets()
    updateEnemies()
    updateObstacles()
    updateParticles()
    checkCollisions()
    updateCamera()
    proceduralSpawn()
    
    if (player.flags & 0x02) == 0 { // not alive
        enterGameOver()
        if score > highScore {
            highScore = score
        }
    }
}

// Responsável pelo funcionamento da munição
func updateAmmo() {
    if isReloading {
        reloadTimer++
        if reloadTimer >= RELOAD_TIME {
            ammo = MAX_AMMO
            isReloading = false
            reloadTimer = 0
        }
    } else if ammo == 0 {
        // Inicia recarga automática quando não há mais munição
        isReloading = true
        reloadTimer = 0
    }
}

func enterGameOver() {
    gameState = STATE_GAME_OVER
    gameOverTimer = 120 // 2 segundos
    previousGamepadState = *GAMEPAD1 // Captura o estado atual dos botões
}

func updateGameOver() {
    gamepad := *GAMEPAD1
    
    // Primeiro, aguarda o timer E que os botões sejam soltos
    if gameOverTimer > 0 {
        gameOverTimer--
        previousGamepadState = gamepad // Atualiza o estado anterior
        return
    }
    
    // Só aceita input se:
    // 1. Algum botão foi pressionado AGORA
    // 2. E NÃO estava pressionado no frame anterior
    buttonJustPressed := (gamepad&(BUTTON_1|BUTTON_2)) != 0 && (previousGamepadState&(BUTTON_1|BUTTON_2)) == 0
    
    if buttonJustPressed {
        gameState = STATE_MENU
        resetGame()
    }
    
    previousGamepadState = gamepad
}

func startGame() {
    resetGame()
    player.flags = 0x03
    gameFrame = 0
    score = 0
    cameraX = 0
    
    // Reset do sistema de munição
    ammo = MAX_AMMO
    reloadTimer = 0
    isReloading = false

    // Reset da direção da mira
    aimDirection = AIM_HORIZONTAL

    // Seed baseado no tempo real absoluto (frameCounter nunca reseta)
    gameStartRealTime = frameCounter
    rngSeed = uint32(gameStartRealTime*31337 + (gameStartRealTime<<7) + (gameStartRealTime>>3))

    // Mistura com outros fatores únicos do momento
    rngSeed = rngSeed ^ uint32(gameStartRealTime*gameStartRealTime)
    rngSeed = rngSeed ^ uint32((gameStartRealTime%997) * 2039)
    rngSeed = rngSeed ^ uint32((gameStartRealTime%1009) * 4093)

    // "Aquece" o gerador com base no tempo
    warmupCount := (gameStartRealTime % 50) + 10
    for i := 0; i < int(warmupCount); i++ {
        rng()
    }
    
    // Reset do sistema procedural com valores mais variados
    spawnTimer = randInt(20) // Começa em momento aleatório
    nextSpawnDelay = 30 + randInt(40) // Delay inicial mais variado
    lastSpawnType = -1
    lastSpawnX = 0
    currentPattern = int8(randInt(4)) // Padrão inicial aleatório
    patternTimer = randInt(400) // Tempo de padrão inicial muito variado
}

func resetGame() {
    initPlayer()
    clearArrays()
}

func handleInput() {
    gamepad := *GAMEPAD1
    mouseButtons := *MOUSE_BUTTONS
    
    // Detectar pressionamento (não repetição)
    gamepadPressed := gamepad & ^prevGamepad
    mousePressed := mouseButtons & ^prevMouseButtons
    
    // Controle da mira - direcionais do gamepad
    if gamepad&BUTTON_RIGHT != 0 {
        aimDirection = AIM_HORIZONTAL
    }
    if gamepad&BUTTON_UP != 0 {
        aimDirection = AIM_VERTICAL
    }
    
    // Pulo - BUTTON_1 (X, V, espaço ou botão esquerdo do mouse)
    if (gamepadPressed&BUTTON_1 != 0 || mousePressed&MOUSE_LEFT != 0) && 
       (player.flags&0x01) != 0 { // onGround
        player.velY = currentJumpPower
        player.flags &= 0xFE // clear onGround
    }
    
    // Tiro - BUTTON_2 (Z, C ou botão direito do mouse)
    if (gamepadPressed&BUTTON_2 != 0 || mousePressed&MOUSE_RIGHT != 0){
        shoot()
    }
    
    // Atualizar estados anteriores
    prevGamepad = gamepad
    prevMouseButtons = mouseButtons
}

// Responsável pela atualização do estado do jogador
func updatePlayer() {
    if (player.flags & 0x02) == 0 { // not alive
        return
    }
    
    // Movimento horizontal - mais rápido no ar para facilitar pulos
    if (player.flags & 0x01) == 0 { // not onGround (está no ar)
        player.x += currentPlayerSpeed
    } else {
        player.x += currentPlayerSpeed
    }
    
    // Gravidade
    if (player.flags & 0x01) == 0 { // not onGround
        player.velY += 1
        if player.velY > 8 {
            player.velY = 8
        }
        player.y += int32(player.velY)
        
        if player.y >= GROUND_Y-PLAYER_HEIGHT {
            player.y = GROUND_Y - PLAYER_HEIGHT
            player.velY = 0
            player.flags |= 0x01 // set onGround
        }
    }
    
    // Animação de corrida
    player.animFrame++
    if player.animFrame > 20 {
        player.animFrame = 0
    }
}

// Responsável pela mecânica de tiro
func updateBullets() {
    for i := 0; i < MAX_BULLETS; i++ {
        if bullets[i].active {
            bullets[i].x += int32(bullets[i].velX)
            bullets[i].y += int32(bullets[i].velY)
            
            // Remove se sair da tela em qualquer direção
            if bullets[i].x < cameraX-20 || bullets[i].x > cameraX+SCREEN_WIDTH+20 ||
               bullets[i].y < -20 || bullets[i].y > SCREEN_HEIGHT+20 {
                bullets[i].active = false
            }
        }
    }
}

// Responsável pela mecânica de tiro inimigo
func updateEnemyBullets() {
    for i := 0; i < MAX_ENEMY_BULLETS; i++ {
        if enemyBullets[i].active {
            enemyBullets[i].x += int32(enemyBullets[i].velX)
            enemyBullets[i].y += int32(enemyBullets[i].velY)
            
            // Remove se sair da tela
            if enemyBullets[i].x < cameraX-20 || enemyBullets[i].x > cameraX+SCREEN_WIDTH+20 ||
               enemyBullets[i].y < 0 || enemyBullets[i].y > SCREEN_HEIGHT {
                enemyBullets[i].active = false
            }
        }
    }
}

// Atualiza as características e estado dos inimigos
func updateEnemies() {
    for i := 0; i < MAX_ENEMIES; i++ {
        if enemies[i].active {
            enemies[i].animFrame++
            
            // Sistema de tiro dos inimigos
            if enemies[i].animFrame%ENEMY_SHOOT_RATE == 0 {
                shootEnemyBullet(enemies[i].x, enemies[i].y, enemies[i].enemyType)
            }

            // Características dos inimigos terrestres
            if enemies[i].enemyType == ENEMY_GROUND {
                enemies[i].velX = -int8(currentEnemySpeed)
                enemies[i].y = GROUND_Y - 12
            } else if enemies[i].enemyType == ENEMY_FLYING {    // Características dos aéreos
                enemies[i].velX = -int8(currentEnemySpeed)
                // Movimento senoidal para voar
                if (enemies[i].animFrame/15)%2 == 0 {
                    enemies[i].velY = -1
                } else {
                    enemies[i].velY = 1
                }
                
                if enemies[i].y < 40 {
                    enemies[i].y = 40
                    enemies[i].velY = 1
                }
                if enemies[i].y > GROUND_Y-30 {
                    enemies[i].y = GROUND_Y - 30
                    enemies[i].velY = -1
                }
            }
            
            enemies[i].x += int32(enemies[i].velX)
            enemies[i].y += int32(enemies[i].velY)
            
            if enemies[i].x < cameraX-50 {
                enemies[i].active = false
            }
        }
    }
}

// Mecanismo de tiro inimigo
func shootEnemyBullet(x, y int32, enemyType int8) {
    for i := 0; i < MAX_ENEMY_BULLETS; i++ {
        if !enemyBullets[i].active {
            if enemyType == ENEMY_GROUND {
                enemyBullets[i].x = x - 2
                enemyBullets[i].y = y + 6
                enemyBullets[i].velX = -2
                enemyBullets[i].velY = 0
            } else {
                enemyBullets[i].x = x + 4
                enemyBullets[i].y = y + 8
                enemyBullets[i].velX = 0
                enemyBullets[i].velY = 1
            }
            
            enemyBullets[i].active = true
            return
        }
    }
}

// Responsável pelos obstáculos
func updateObstacles() {
    for i := 0; i < MAX_OBSTACLES; i++ {
        if obstacles[i].active {
            if obstacles[i].x < cameraX-50 {
                obstacles[i].active = false
            }
        }
    }
}

// Responsável pelas partículas
func updateParticles() {
    for i := 0; i < MAX_PARTICLES; i++ {
        if particles[i].active {
            particles[i].x += int32(particles[i].velX)
            particles[i].y += int32(particles[i].velY)
            particles[i].velY += 1
            particles[i].life--
            
            if particles[i].life <= 0 {
                particles[i].active = false
            }
        }
    }
}

// Verifica as diversas colisões possíveis
func checkCollisions() {
    // Tiro vs inimigo
    for i := 0; i < MAX_BULLETS; i++ {
        if bullets[i].active {
            for j := 0; j < MAX_ENEMIES; j++ {
                if enemies[j].active {
                    var enemyW, enemyH int32 = 8, 12
                    if enemies[j].enemyType == ENEMY_FLYING {
                        enemyH = 8
                    }
                    
                    if collision(bullets[i].x, bullets[i].y, BULLET_WIDTH, BULLET_HEIGHT,
                                enemies[j].x, enemies[j].y, enemyW, enemyH) {
                        bullets[i].active = false
                        enemies[j].active = false
                        score += 10
                        createExplosion(enemies[j].x, enemies[j].y)
                    }
                }
            }
        }
    }
    
    // Tiro vs tiro inimigo
    for i := 0; i < MAX_BULLETS; i++ {
        if bullets[i].active {
            for j := 0; j < MAX_ENEMY_BULLETS; j++ {
                if enemyBullets[j].active {
                    if collision(bullets[i].x-1, bullets[i].y-1, BULLET_WIDTH+2, BULLET_HEIGHT+2,
                                enemyBullets[j].x-1, enemyBullets[j].y-1, ENEMY_BULLET_WIDTH+2, ENEMY_BULLET_HEIGHT+2) {
                        bullets[i].active = false
                        enemyBullets[j].active = false
                        createExplosion(bullets[i].x, bullets[i].y)
                        score += 5 // Bônus por interceptar bala inimiga
                    }
                }
            }
        }
    }
    
    // Jogador vs tiro inimigo
    for i := 0; i < MAX_ENEMY_BULLETS; i++ {
        if enemyBullets[i].active {
            if collision(player.x, player.y, PLAYER_WIDTH, PLAYER_HEIGHT,
                        enemyBullets[i].x, enemyBullets[i].y, ENEMY_BULLET_WIDTH, ENEMY_BULLET_HEIGHT) {
                enemyBullets[i].active = false
                player.flags &= 0xFD // clear alive
                createExplosion(player.x, player.y)
            }
        }
    }
    
    // Jogador vs inimigos
    for i := 0; i < MAX_ENEMIES; i++ {
        if enemies[i].active {
            var enemyW, enemyH int32 = 8, 12
            if enemies[i].enemyType == ENEMY_FLYING {
                enemyH = 8
            }
            
            if collision(player.x, player.y, PLAYER_WIDTH, PLAYER_HEIGHT,
                        enemies[i].x, enemies[i].y, enemyW, enemyH) {
                player.flags &= 0xFD // clear alive
                createExplosion(player.x, player.y)
            }
        }
    }
    
    // Jogador vs obstáculos
    for i := 0; i < MAX_OBSTACLES; i++ {
        if obstacles[i].active {
            if collision(player.x+1, player.y+2, PLAYER_WIDTH-2, PLAYER_HEIGHT-2,
                        obstacles[i].x, obstacles[i].y, int32(obstacles[i].width), int32(obstacles[i].height)) {
                player.flags &= 0xFD // clear alive
                createExplosion(player.x, player.y)
            }
        }
    }
}

// Atualiza a camera
func updateCamera() {
    cameraX = player.x - SCREEN_WIDTH/8
    if cameraX < 0 {
        cameraX = 0
    }
}

// Auxiliar da check colision
func collision(x1, y1, w1, h1, x2, y2, w2, h2 int32) bool {
    return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}

// Mecanismo de tiro do jogador
func shoot() {
    if ammo <= 0 || isReloading {
        return // Não pode atirar se não tem munição ou está recarregando
    }
    
    for i := 0; i < MAX_BULLETS; i++ {
        if !bullets[i].active {
            if aimDirection == AIM_HORIZONTAL {
                // Tiro horizontal
                bullets[i].x = player.x + PLAYER_WIDTH
                bullets[i].y = player.y + 4
                bullets[i].velX = 5
                bullets[i].velY = 0
            } else {
                // Tiro vertical
                bullets[i].x = player.x + 4 // centralizado no player
                bullets[i].y = player.y - 2 // um pouco acima
                bullets[i].velX = 0
                bullets[i].velY = -5 // velocidade para cima
            }
            bullets[i].active = true
            ammo-- // Consome munição
            return
        }
    }
}

// Geração de inimigos
func spawnEnemy(x, y int32, enemyType int8) {
    for i := 0; i < MAX_ENEMIES; i++ {
        if !enemies[i].active {
            enemies[i].x = x
            enemies[i].y = y
            enemies[i].enemyType = enemyType
            enemies[i].active = true
            enemies[i].animFrame = 0
            return
        }
    }
}

// Geração de obstáculos
func spawnObstacle(x, y int32, width, height int8, obstacleType int8) {
    for i := 0; i < MAX_OBSTACLES; i++ {
        if !obstacles[i].active {
            obstacles[i].x = x
            obstacles[i].y = y
            obstacles[i].width = width
            obstacles[i].height = height
            obstacles[i].obstacleType = obstacleType
            obstacles[i].active = true
            return
        }
    }
}

func createExplosion(x, y int32) {
    for i := 0; i < MAX_PARTICLES; i++ {
        if !particles[i].active {
            particles[i].x = x + int32(i*2)
            particles[i].y = y + int32(i)
            particles[i].velX = int8((i%3 - 1) * 2)
            particles[i].velY = int8(-2 - i/2)
            particles[i].life = 25
            particles[i].active = true
            return
        }
    }
}

func draw() {
    switch gameState {
    case STATE_MENU:
        drawMenu()
    case STATE_PLAYING:
        drawGame()
    case STATE_GAME_OVER:
        drawGameOver()
    }
}

func drawMenu() {
    *DRAW_COLORS = 0x01
    rect(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT)
    
    *DRAW_COLORS = 0x03
    drawSimpleText("JUMP 'N' SHOOT", 40, 30)
    
    *DRAW_COLORS = 0x04
    drawSimpleText("PRESS ANY BUTTON", 35, 60)
    
    *DRAW_COLORS = 0x04
    drawSimpleText("HIGH:", 60, 90)
    drawNumber(highScore, 100, 90)

    *DRAW_COLORS = 0x03
    drawSimpleText("JUMP:X,V,SPACE,MOUSE(LEFT)", 5, 120)
    drawSimpleText("SHOOT:Z,C,MOUSE(RIGHT)", 5, 130)
}

func drawGame() {
    // Céu
    *DRAW_COLORS = 0x01
    rect(0, 0, SCREEN_WIDTH, GROUND_Y-20)
    
    *DRAW_COLORS = 0x10
    rect(0, GROUND_Y-20, SCREEN_WIDTH, 20)
    
    // Chão
    *DRAW_COLORS = 0x02
    rect(0, GROUND_Y, SCREEN_WIDTH, SCREEN_HEIGHT-GROUND_Y)
    
    drawPlayer()
    drawBullets()
    drawEnemyBullets()
    drawEnemies()
    drawObstacles()
    drawParticles()
    drawUI()
}

func drawGameOver() {
    *DRAW_COLORS = 0x01
    rect(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT)
    
    *DRAW_COLORS = 0x03
    drawSimpleText("GAME OVER", 50, 40)
    drawSimpleText("SCORE:", 50, 80)
    drawNumber(score, 110, 80)
    drawSimpleText("PRESS ANY BUTTON", 35, 120)
}

func drawPlayer() {
    if (player.flags & 0x03) == 0 { // not alive
        return
    }
    
    screenX := player.x - cameraX
    if screenX < -PLAYER_WIDTH || screenX > SCREEN_WIDTH {
        return
    }
    
    // Selecionar frame de animação
    frame := 0
    if player.animFrame > 10 {
        frame = 1
    }
    
    // Desenhar player
    drawSprite8x12(playerSprite[frame][:], screenX, player.y, 0x03)
    
    // Desenhar arma na posição adequada
    if aimDirection == AIM_HORIZONTAL {
        drawWeapon(screenX + 8, player.y + 4)
    } else {
        drawWeapon(screenX + 3, player.y + 4) // mais centralizado para arma vertical
    }
}

func drawBulletIcon(x, y int32) {
    *DRAW_COLORS = 0x04
    for row := 0; row < 2; row++ {
        data := bulletSprite[row]
        for col := 0; col < 4; col++ {
            if (data & (1 << (7 - col))) != 0 {
                rect(x+int32(col), y+int32(row), 1, 1)
            }
        }
    }
}

func drawWeapon(x, y int32) {
    *DRAW_COLORS = 0x3
    
    if aimDirection == AIM_HORIZONTAL {
        // Arma horizontal
        for row := 0; row < 3; row++ {
            data := weaponSprite[row]
            for col := 0; col < 6; col++ {
                if (data & (1 << (7 - col))) != 0 {
                    rect(x+int32(col), y+int32(row), 1, 1)
                }
            }
        }
    } else {
        // Arma vertical
        for row := 0; row < 6; row++ {
            data := weaponVerticalSprite[row]
            for col := 0; col < 3; col++ {
                if (data & (1 << (7 - col))) != 0 {
                    rect(x+int32(col), y-6+int32(row), 1, 1) // ajusta posição para cima
                }
            }
        }
    }
}

func drawBullets() {
    for i := 0; i < MAX_BULLETS; i++ {
        if bullets[i].active {
            screenX := bullets[i].x - cameraX
            if screenX >= -10 && screenX < SCREEN_WIDTH+10 {
                if bullets[i].velY != 0 {
                    // Bala vertical
                    *DRAW_COLORS = 0x04
                    rect(screenX, bullets[i].y, BULLET_HEIGHT, BULLET_WIDTH) // invertido
                    *DRAW_COLORS = 0x03
                    rect(screenX, bullets[i].y, BULLET_HEIGHT, 2) // ponta
                } else {
                    // Bala horizontal
                    *DRAW_COLORS = 0x04
                    rect(screenX, bullets[i].y, BULLET_WIDTH, BULLET_HEIGHT)
                    *DRAW_COLORS = 0x03
                    rect(screenX + BULLET_WIDTH - 2, bullets[i].y, 2, BULLET_HEIGHT)
                }
            }
        }
    }
}

func drawEnemies() {
    for i := 0; i < MAX_ENEMIES; i++ {
        if enemies[i].active {
            screenX := enemies[i].x - cameraX
            if screenX >= -20 && screenX < SCREEN_WIDTH+20 {
                if enemies[i].enemyType == ENEMY_GROUND {
                    // Inimigo terrestre
                    drawSprite8x12(groundEnemySprite[:], screenX, enemies[i].y, 0x03)
                } else {
                    // Inimigo voador
                    drawSprite8x8(flyingEnemySprite[:], screenX, enemies[i].y, 0x03)
                }
            }
        }
    }
}

func drawEnemyBullets() {
    for i := 0; i < MAX_ENEMY_BULLETS; i++ {
        if enemyBullets[i].active {
            screenX := enemyBullets[i].x - cameraX
            if screenX >= -10 && screenX < SCREEN_WIDTH+10 {
                // Balas dos inimigos
                *DRAW_COLORS = 0x03
                rect(screenX, enemyBullets[i].y, ENEMY_BULLET_WIDTH, ENEMY_BULLET_HEIGHT)
                // Adicionar um pixel central mais brilhante para melhor visibilidade
                *DRAW_COLORS = 0x04
                rect(screenX+1, enemyBullets[i].y+1, 1, 1)
            }
        }
    }
}

func drawObstacles() {
    for i := 0; i < MAX_OBSTACLES; i++ {
        if obstacles[i].active {
            screenX := obstacles[i].x - cameraX
            if screenX >= -30 && screenX < SCREEN_WIDTH+30 {
                if obstacles[i].obstacleType == OBSTACLE_ROCK {
                    // Rocha
                    drawSprite8x8(rockSprite[:], screenX, obstacles[i].y, 0x04)
                } else {
                    // Spike
                    drawSprite6x8(spikeSprite[:], screenX, obstacles[i].y, 0x04)
                }
            }
        }
    }
}

func drawParticles() {
    for i := 0; i < MAX_PARTICLES; i++ {
        if particles[i].active {
            screenX := particles[i].x - cameraX
            if screenX >= 0 && screenX < SCREEN_WIDTH {
                *DRAW_COLORS = 0x32
                rect(screenX, particles[i].y, 2, 2)
                *DRAW_COLORS = 0x03
                rect(screenX, particles[i].y, 1, 1)
            }
        }
    }
}

func drawUI() {
    *DRAW_COLORS = 0x03
    drawSimpleText("SCORE:", 5, 5)
    drawNumber(score, 45, 5)
    
    // Indicador de direção da mira
    *DRAW_COLORS = 0x04
    if aimDirection == AIM_HORIZONTAL {
        drawSimpleText("AIM: FORWARD", 80, 5)
    } else {
        drawSimpleText("AIM: UP", 80, 5)
    }
    
    // Indicador de munição
    *DRAW_COLORS = 0x04
    drawSimpleText("AMMO:", 5, 15)
    
    if isReloading {
        *DRAW_COLORS = 0x03
        drawSimpleText("RELOAD", 45, 15)
        // Barra de progresso do reload
        *DRAW_COLORS = 0x02
        rect(45, 25, 60, 4)
        *DRAW_COLORS = 0x04
        progress := (reloadTimer * 60) / RELOAD_TIME
        rect(45, 25, progress, 4)
    } else {
        // Desenhar balas restantes
        for i := 0; i < int(ammo); i++ {
            drawBulletIcon(45+int32(i*6), 15)
        }
    }
}

// Função para desenhar sprite 8x8
func drawSprite8x8(sprite []uint8, x, y int32, colors uint16) {
    *DRAW_COLORS = colors
    for row := 0; row < 8; row++ {
        data := sprite[row]
        for col := 0; col < 8; col++ {
            if (data & (1 << (7 - col))) != 0 {
                rect(x+int32(col), y+int32(row), 1, 1)
            }
        }
    }
}

// Função para desenhar sprite 6x8
func drawSprite6x8(sprite []uint8, x, y int32, colors uint16) {
    *DRAW_COLORS = colors
    for row := 0; row < 8; row++ {
        data := sprite[row]
        for col := 0; col < 6; col++ {
            if (data & (1 << (7 - col))) != 0 {
                rect(x+int32(col), y+int32(row), 1, 1)
            }
        }
    }
}

// Função para desenhar sprite 8x12
func drawSprite8x12(sprite []uint8, x, y int32, colors uint16) {
    *DRAW_COLORS = colors
    for row := 0; row < 12; row++ {
        data := sprite[row]
        for col := 0; col < 8; col++ {
            if (data & (1 << (7 - col))) != 0 {
                rect(x+int32(col), y+int32(row), 1, 1)
            }
        }
    }
}

func drawSimpleText(text string, x, y int32) {
    for i := 0; i < len(text); i++ {
        drawSimpleChar(text[i], x+int32(i)*6, y)
    }
}

func drawSimpleChar(char byte, x, y int32) {
    switch char {
    case 'A':
        rect(x, y+1, 1, 4)
        rect(x+1, y, 2, 1)
        rect(x+1, y+2, 2, 1)
        rect(x+3, y+1, 1, 4)
    case 'B':
        rect(x, y, 1, 5)
        rect(x+1, y, 2, 1)
        rect(x+1, y+2, 2, 1)
        rect(x+1, y+4, 2, 1)
        rect(x+3, y+1, 1, 1)
        rect(x+3, y+3, 1, 1)
    case 'C':
        rect(x, y+1, 1, 3)
        rect(x+1, y, 2, 1)
        rect(x+1, y+4, 2, 1)
    case 'D':
        rect(x, y, 1, 5)
        rect(x+1, y, 2, 1)
        rect(x+1, y+4, 2, 1)
        rect(x+3, y+1, 1, 3)
    case 'E':
        rect(x, y, 1, 5)
        rect(x+1, y, 2, 1)
        rect(x+1, y+2, 2, 1)
        rect(x+1, y+4, 2, 1)
    case 'F':
        rect(x, y, 1, 5)
        rect(x+1, y, 2, 1)
        rect(x+1, y+2, 2, 1)
    case 'G':
        rect(x, y+1, 1, 3)
        rect(x+1, y, 2, 1)
        rect(x+1, y+4, 2, 1)
        rect(x+2, y+2, 2, 1)
        rect(x+3, y+2, 1, 2)
    case 'H':
        rect(x, y, 1, 5)
        rect(x+1, y+2, 2, 1)
        rect(x+3, y, 1, 5)
    case 'I':
        rect(x, y, 3, 1)
        rect(x+1, y+1, 1, 3)
        rect(x, y+4, 3, 1)
    case 'J':
        rect(x+2, y, 1, 4)
        rect(x, y+4, 3, 1)
        rect(x, y+2, 1, 2)
    case 'K':
        rect(x, y, 1, 5)
        rect(x+1, y+2, 1, 1)
        rect(x+2, y+1, 1, 1)
        rect(x+2, y+3, 1, 1)
        rect(x+3, y, 1, 1)
        rect(x+3, y+4, 1, 1)
    case 'L':
        rect(x, y, 1, 5)
        rect(x+1, y+4, 2, 1)
    case 'M':
        rect(x, y, 1, 5)
        rect(x+1, y+1, 1, 1)
        rect(x+2, y+2, 1, 1)
        rect(x+3, y+1, 1, 1)
        rect(x+4, y, 1, 5)
    case 'N':
        rect(x, y, 1, 5)
        rect(x+1, y+1, 1, 1)
        rect(x+2, y+2, 1, 1)
        rect(x+3, y, 1, 5)
    case 'O':
        rect(x, y+1, 1, 3)
        rect(x+1, y, 2, 1)
        rect(x+1, y+4, 2, 1)
        rect(x+3, y+1, 1, 3)
    case 'P':
        rect(x, y, 1, 5)
        rect(x+1, y, 2, 1)
        rect(x+1, y+2, 2, 1)
        rect(x+3, y+1, 1, 1)
    case 'Q':
        rect(x, y+1, 1, 3)
        rect(x+1, y, 2, 1)
        rect(x+1, y+4, 2, 1)
        rect(x+3, y+1, 1, 2)
        rect(x+2, y+3, 1, 1)
        rect(x+3, y+4, 1, 1)
    case 'R':
        rect(x, y, 1, 5)      
        rect(x+1, y, 2, 1)    
        rect(x+1, y+2, 2, 1)  
        rect(x+3, y+1, 1, 1)  
        rect(x+2, y+3, 1, 1) 
        rect(x+3, y+4, 1, 1)  
    case 'S':
        rect(x, y, 3, 1)
        rect(x, y+1, 1, 1)
        rect(x, y+2, 3, 1)
        rect(x+2, y+3, 1, 1)
        rect(x, y+4, 3, 1)
    case 'T':
        rect(x, y, 3, 1)
        rect(x+1, y+1, 1, 4)
    case 'U':
        rect(x, y, 1, 4)
        rect(x+1, y+4, 2, 1)
        rect(x+3, y, 1, 4)
    case 'V':
        rect(x, y, 1, 3)
        rect(x+1, y+3, 1, 1)
        rect(x+2, y+4, 1, 1)
        rect(x+3, y+3, 1, 1)
        rect(x+4, y, 1, 3)
    case 'W':
        rect(x, y, 1, 5)
        rect(x+1, y+4, 1, 1)
        rect(x+2, y+3, 1, 1)
        rect(x+3, y+4, 1, 1)
        rect(x+4, y, 1, 5)
    case 'X':
        rect(x, y, 1, 1)
        rect(x+1, y+1, 1, 1)
        rect(x+2, y+2, 1, 1)
        rect(x+3, y+1, 1, 1)
        rect(x+4, y, 1, 1)
        rect(x, y+4, 1, 1)
        rect(x+1, y+3, 1, 1)
        rect(x+3, y+3, 1, 1)
        rect(x+4, y+4, 1, 1)
    case 'Y':
        rect(x, y, 1, 1)
        rect(x+1, y+1, 1, 1)
        rect(x+2, y+2, 1, 3)
        rect(x+3, y+1, 1, 1)
        rect(x+4, y, 1, 1)
    case 'Z':
        rect(x, y, 4, 1)
        rect(x+3, y+1, 1, 1)
        rect(x+2, y+2, 1, 1)
        rect(x+1, y+3, 1, 1)
        rect(x, y+4, 4, 1)
    case ' ':
        // espaço em branco
    case ':':
        rect(x+1, y+1, 1, 1)
        rect(x+1, y+3, 1, 1)
    }
}

func drawNumber(num, x, y int32) {
    if num == 0 {
        drawDigit(0, x, y)
        return
    }
    
    if num >= 10 {
        drawNumber(num/10, x-6, y)
    }
    drawDigit(int(num%10), x, y)
}

func drawDigit(digit int, x, y int32) {
    switch digit {
    case 0:
        rect(x, y, 3, 1)
        rect(x, y+1, 1, 3)
        rect(x+2, y+1, 1, 3)
        rect(x, y+4, 3, 1)
    case 1:
        rect(x+1, y, 1, 5)
    case 2:
        rect(x, y, 3, 1)
        rect(x+2, y+1, 1, 1)
        rect(x, y+2, 3, 1)
        rect(x, y+3, 1, 1)
        rect(x, y+4, 3, 1)
    case 3:
        rect(x, y, 3, 1)
        rect(x+2, y+1, 1, 1)
        rect(x+1, y+2, 2, 1)
        rect(x+2, y+3, 1, 1)
        rect(x, y+4, 3, 1)
    case 4:
        rect(x, y, 1, 3)
        rect(x+2, y, 1, 5)
        rect(x+1, y+2, 1, 1)
    case 5:
        rect(x, y, 3, 1)
        rect(x, y+1, 1, 1)
        rect(x, y+2, 3, 1)
        rect(x+2, y+3, 1, 1)
        rect(x, y+4, 3, 1)
    case 6:
        rect(x, y, 3, 1)
        rect(x, y+1, 1, 1)
        rect(x, y+2, 3, 1)
        rect(x, y+3, 1, 1)
        rect(x+2, y+3, 1, 1)
        rect(x, y+4, 3, 1)
    case 7:
        rect(x, y, 3, 1)
        rect(x+2, y+1, 1, 4)
    case 8:
        rect(x, y, 3, 1)
        rect(x, y+1, 1, 1)
        rect(x+2, y+1, 1, 1)
        rect(x, y+2, 3, 1)
        rect(x, y+3, 1, 1)
        rect(x+2, y+3, 1, 1)
        rect(x, y+4, 3, 1)
    case 9:
        rect(x, y, 3, 1)
        rect(x, y+1, 1, 1)
        rect(x+2, y+1, 1, 1)
        rect(x, y+2, 3, 1)
        rect(x+2, y+3, 1, 1)
        rect(x, y+4, 3, 1)
    }
}

//go:wasmimport env rect
func rect(x, y, width, height int32)

func main() {}
