## 合约设计

### 规则说明

编写solidity合约，实现一个游戏Demo：
- 条件

    - 玩家数量：2
    - 游戏道具：游戏开始时，每个玩家持有卡牌：随机获取3个（剪刀/岩石/布）
- 游戏规则：
    - 回合开始：双方玩家需要从他们的手中决定使用哪个卡牌（剪刀/岩石/布），而不让其他玩家知道卡牌代表什么；
    - 双方都决定使用什么卡牌后，系统会揭示他们使用卡牌的内容
    - 根据以下规则结算：剪刀>布；石头>剪刀；布>石头
    - 结算后，输掉的一方需要丢弃一个生命点卡牌
    - 游戏结束后，那一轮使用的卡牌也需要被丢弃并重复上述步骤
    - 直到一方用完所有生命点卡牌，就被认为是另一方的胜利

### 流程图

![][游戏demo流程图]

### 主要参数定义
- 定义卡牌类型
```js
enum CardType {
    None,   // 0: 空白
    Rock,   // 1：石头
    Scissors, // 2：剪刀
    Paper  // 3：布
}
```

- 定义游戏桌状态
```js
enum TableState {
    PreReady,  // 0: 预准备状态
    Ready,     // 1：准备状态
    Active,    // 2: 游戏中
    Over        // 3:游戏结束
}
```

- 定义玩家信息
```js
struct Player {
    bool isInTable; // 是否在游戏桌
    CardType[] cards; // 卡牌
}
```

- 定义游戏桌信息
```js
struct Table {
    uint256 id;     // table id
    uint8 round;  // 当前轮次，初始化为0，从1开始
    CardType[][] battleCards; // 对战牌，包括每一轮，第一张为玩家1的牌，第二张为玩家2的牌
    // 玩家的地址列表
    address[] playerAddrList;
    // 保存玩家地址的映射关系
    mapping(address => Player) players;
    TableState state;   // 游戏桌状态
}
```

### 接口说明

#### 创建游戏桌

```js
function createTable() public
```

#### 查询游戏桌详情

```js
function getTableInfo(uint256 tableId) public view 
    returns (
        TableState state, // 游戏桌的状态
        uint8 round,  // 当前轮次，从0开始，不能超过cardNums
        CardType[][] memory battleCards,// 对战牌，包括每一轮
        address[] memory playerAddrList // 玩家的地址列表
    ) 
```
> 参数说明：
>
> - tableId：游戏桌Id

#### 查询游戏桌上的玩家信息

```js
function getPlayerInfo(uint256 tableId, address player) public view 
    returns (
        bool isInTable,
        CardType[] memory cards // 卡牌
    ) {
    Table storage table = tables[tableId];
    isInTable = table.players[player].isInTable;
    cards = table.players[player].cards;
}
```
> 参数说明：
>
> - tableId：游戏桌Id；
> - player：玩家地址

#### 玩家加入游戏桌

```js
function joinTable(uint256 tableId) external
```
> 参数说明：
>
> - tableId：游戏桌ID；

#### 玩家退出游戏桌

```js
function exitTable(uint256 tableId) external
```
> 参数说明：
>
> - tableId：游戏桌ID；

#### 玩家准备就绪动作

注意：只能是玩家2进行操作；

```js
function readyGame(uint256 tableId) external
```
> 参数说明：
>
> - tableId：游戏桌ID；

#### 开始游戏

```js
function startGame(uint256 tableId) external
```
注意：只能是玩家1进行操作，并且必须等到玩家2操作readyGame成功之后才能进行操作；

> 参数说明：
>
> - tableId：游戏桌ID；

#### 玩家出牌
说明：当每轮游戏两个玩家都出牌后，自动进行当前轮的比较和执行卡牌的处理逻辑。

```js
function playCard(uint256 tableId, uint256 index) external
```
> 参数说明：
>
> - tableId：游戏桌ID；
> - index：玩家卡牌列表中的索引下标；


## 合约事件

```js
/// @notice 创建游戏桌事件
/// @dev 创建游戏桌事件
/// @param _tableId 游戏桌ID
event CreateTable(uint256 _tableId);

/// @notice 玩家加入游戏桌
/// @dev 玩家加入游戏桌，当第二个人加入时，自动启动游戏
/// @param _tableId 游戏桌ID
/// @param _player 玩家地址
event JoinTable(uint256 _tableId, address _player);

/// @notice 玩家退出游戏桌
/// @dev 玩家退出游戏桌，游戏未开始或结束时，可退出
/// @param _tableId 游戏桌ID
/// @param _player 玩家地址
event ExitTable(uint256 _tableId, address _player);

/// @notice 新启动一局游戏的事件
/// @dev 新启动一局游戏的事件
/// @param _tableId 游戏桌ID
/// @param _playerAddrList 玩家地址列表
event StartRound(uint256 _tableId, address[] _playerAddrList);

/// @notice 结束游戏轮触发事件
/// @dev 调用endRound接口时触发
/// @param _tableId 游戏桌ID
/// @param _round 当前轮
/// @param _winner 赢家
/// @param _cards 当前轮对战牌
event EndRound(uint256 indexed _tableId, uint8 indexed _round, address _winner,     
  CardType[2] _cards);

/// @notice 游戏结束事件
/// @dev 重置游戏相关的参数
/// @param _tableId 游戏桌ID
/// @param _winner 赢家地址,平局时，地址为0x0
event GameOver(uint256 _tableId, address _winner);
```



[游戏demo流程图]: ./resource/RockPaperScissorsGame_flowchar.png

