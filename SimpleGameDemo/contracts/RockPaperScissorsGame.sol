// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "./BehalfSignature.sol";
/**
此游戏合约遵循如下规则：
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
*/

contract RockPaperScissorsGame is BehalfSignature{
    // 初始化几桌游戏桌
    uint private constant initNum = 2;
    // 每个玩家每局拥有的牌数
    uint8 private constant cardNums = 3;

    // 定义牌类型
    enum CardType {
        None,   // 0
        Rock,   // 1 
        Scissors, // 2
        Paper  // 3
    }

    // 游戏桌状态
    enum TableState {
        PreReady,  // 0: 预准备状态
        Ready,     // 1：准备状态
        Active,    // 2: 游戏中
        Over        // 3:游戏结束
    }

    // 定义玩家信息
    struct Player {
        bool isInTable; // 是否在游戏桌
        CardType[] cards; // 卡牌
    }


    // 定义游戏桌结构体
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
    // 玩家信息
    // mapping(address => Player) public players;
    // 定义游戏桌数组
    mapping (uint256 => Table) public tables;
    // 桌号（从0开始）
    uint256 public numTable = 0;

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

    /// @notice 初始化玩家信息
    /// @dev 初始化玩家信息
    /// @param _tableId 游戏桌ID
    /// @param _player 玩家地址
    event InitPlayer(uint256 _tableId, address _player);

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
    event EndRound(uint256 indexed _tableId, uint8 indexed _round, address _winner, CardType[2] _cards);

    /// @notice 游戏结束事件
    /// @dev 重置游戏相关的参数
    /// @param _tableId 游戏桌ID
    /// @param _winner 赢家地址,平局时，地址为0x0
    event GameOver(uint256 _tableId, address _winner);
    
    constructor() {
        // 初始化游戏
        // 创建多张游戏桌
        for(uint i = 0; i < initNum; i++) {
            createTable();
        }
    }

    // 获取随机数
    function getRandomNumber(uint256 num) internal view returns (uint256[] memory randomList) {
        require(block.number > num, "Getting too many random numbers");
        randomList = new uint256[](num);
        for(uint i = 0; i < num; i++) { 
            uint256 blockValue = uint256(blockhash(block.number - 1 - i)); // 获取前一个区块的哈希值
            uint256 randomNumber = blockValue % 3 + 1; // 获取取随机数
            randomList[i] = randomNumber;
        }
        
        return randomList;
    }

    // 创建游戏桌
    function createTable() public {
        Table storage table = tables[numTable];
        table.id = numTable;
        numTable++;
        emit CreateTable(table.id);
    }

    modifier isExistTable(uint256 tableId) {
        require(tableId < numTable, "Table is not Exist");
        _;
    }
    
    // 加入游戏桌
    function joinTable(uint256 tableId) isExistTable(tableId) external {
        Table storage table = tables[tableId];
        require(!table.players[msg.sender].isInTable, "You are already a player at the game table");
        require(table.playerAddrList.length < 2, "The number of players has exceeded the maximum limit.");
        // 保存玩家信息
        Player memory newPlayer;
        newPlayer.isInTable = true;
        table.players[msg.sender] = newPlayer;
        table.playerAddrList.push(msg.sender);
        emit JoinTable(tableId, msg.sender);
    }

    // 退出游戏桌
    function exitTable(uint256 tableId) isExistTable(tableId) external {
        Table storage table = tables[tableId];
        require(table.players[msg.sender].isInTable, "The player is not at the game table");
        require(table.state == TableState.PreReady || table.state == TableState.Over, "The game is in progress, do not leave the game table");

        // 删除桌面上的玩家地址
        delete table.players[msg.sender];

        // 删除玩家信息
        uint256 index = 0;  // 被删除的位置
        uint256 popIndex = table.playerAddrList.length - 1;
        for(uint256 i = 0; i <= popIndex; i++) {
            if(table.playerAddrList[i] == msg.sender) {
                index = i;
                break;
            }
        }

        if(index < popIndex) {
            table.playerAddrList[index] = table.playerAddrList[popIndex];
        }
        
        table.playerAddrList.pop();
        emit ExitTable(tableId, msg.sender);
    }

    // 初始化玩家信息（牌和生命值）
    function initPlayer(Table storage table, address player) internal {
        uint256 cardNum = 3;
        uint256[] memory cards = getRandomNumber(cardNum);
        for(uint i = 0; i < cardNum; i++ ) {
            table.players[player].cards.push(CardType(cards[i]));
        }
        // table.players[player].cards.push(CardType.Rock);
        // table.players[player].cards.push(CardType.Scissors);
        // table.players[player].cards.push(CardType.Paper);

        emit InitPlayer(table.id, player);
    }

    // 获取游戏桌信息
    function getTableInfo(uint256 tableId) isExistTable(tableId) public view 
        returns (
            TableState state, // 游戏桌的状态
            uint8 round,  // 当前轮次，从0开始，不能超过cardNums
            CardType[][] memory battleCards,// 对战牌，包括每一轮
            address[] memory playerAddrList // 玩家的地址列表
        ) 
    {
        Table storage table = tables[tableId];
        state = table.state;
        round = table.round;
        battleCards = table.battleCards;
        playerAddrList = table.playerAddrList;
    }

    // 获取玩家信息
    function getPlayerInfo(uint256 tableId, address player) isExistTable(tableId) public view 
        returns (
            bool isInTable,
            CardType[] memory cards // 卡牌
        ) {
        Table storage table = tables[tableId];
        isInTable = table.players[player].isInTable;
        cards = table.players[player].cards;
    }

    // 移除玩家的牌
    function removeCard(Table storage table, address player, uint256 index) internal {
        uint256 cardNum = table.players[player].cards.length;
        require(cardNum > index, "removeCard:index error");
        // 将需要清除的元素替换为数组中的最后一个元素
        if(index != cardNum-1) {
            for(uint256 i = index; i < cardNum - 1; i++) {
                table.players[player].cards[i] = table.players[player].cards[i+1];
            }
        }

        // 缩小数组的长度，删除最后一个元素
        table.players[player].cards.pop();
    }

    // 重置游戏桌上的信息
    function resetTable(Table storage table) internal {
        require(table.state == TableState.Over, "The game is not over, can not clear the table information");
        table.round = 0;
        table.state = TableState.PreReady;
        // 重置对战信息
        CardType[][] memory emptyArray = new CardType[][](0);
        table.battleCards = emptyArray;
        // 重置玩家信息
        for (uint256 i = 0; i < table.playerAddrList.length; i++) {
            address addr = table.playerAddrList[i];
            // 清空玩家手里的卡牌
            CardType[] memory emptyCards = new CardType[](0);
            table.players[addr].cards = emptyCards;
        }
    }

    // 游戏准备就绪（后加入桌子的玩家执行）
    function readyGame(uint256 tableId) isExistTable(tableId) external {
        Table storage table = tables[tableId];
        require(table.playerAddrList.length == 2, "There are not enough people at the game table to get ready");
        require(table.state == TableState.PreReady, "The game table is not in a pre-ready state. Ready cannot be performed");
        require(msg.sender == table.playerAddrList[1], "Non-party address cannot execute game ready");
        // 初始化
        initPlayer(table, msg.sender);
        // 修改游戏桌状态为Ready
        table.state = TableState.Ready;
    }

    // 开始游戏（先加入桌子的玩家执行）
    function startGame(uint256 tableId) isExistTable(tableId) external {
        Table storage table = tables[tableId];
        require(table.state == TableState.Ready, "The game table is not ready.");
        require(table.playerAddrList.length == 2, "Not enough players to start the game.");
        require(table.playerAddrList[0] == msg.sender, "Players who do not join the game table first cannot perform the start game.");
        // 初始化
        initPlayer(table, msg.sender);
        // 修改游戏桌状态为Active
        table.state = TableState.Active;
        // 添加第一行
        table.battleCards.push();
        emit StartRound(tableId, table.playerAddrList);
    }

    // 玩家出牌, index表示牌的索引下标，从0开始
    function playCard(uint256 tableId, uint256 index) isExistTable(tableId) external {
        Table storage table = tables[tableId];
        uint8 round = table.round;
        require(table.state == TableState.Active, "The game table is not active.");
        // 取出用户手牌
        Player memory player = table.players[msg.sender];
        require(player.isInTable, "You're not at the game table,You can't play");
        // 索引在玩家手里的牌的范围内
        require(player.cards.length > index, "The specified card index is wrong, out of the range of the player's card");
        // 对战牌位置已经有卡牌，即玩家已经出牌
        uint8 playerIndex = 0;
        if(table.playerAddrList[1] == msg.sender) {
            playerIndex = 1;
        }
        // 判断玩家在当前轮是否出牌
        CardType[] storage cards = table.battleCards[round];
        if(0 == table.battleCards[round].length) {
            // 新增两列
            cards.push(CardType.None);
            cards.push(CardType.None);
        }
        require(CardType.None == cards[playerIndex], "The player has played in the current round");
        
        // 保存用户的牌到桌上
        cards[playerIndex] = player.cards[index];
        // 先移除玩家使用的牌
        removeCard(table, msg.sender, index);

        if(CardType.None != cards[0] && CardType.None != cards[1]) {
            // 对战
            playRound(table, round);
            // 游戏结束
            if(table.state == TableState.Over) {
                resetTable(table);
            } else {
                // 添加一行
                table.battleCards.push();
                // 下一轮(等结算完成)
                table.round++;
            }
        }
    }

    // 两张牌对战
    function playRound(Table storage table, uint8 round) internal {

        // 结算当前轮
        CardType card1 = table.battleCards[round][0];  // 玩家1的牌
        CardType card2 = table.battleCards[round][1];  // 玩家2的牌
        address player1 = table.playerAddrList[0];
        address player2 = table.playerAddrList[1];
        
        // 比较当前轮：当前轮赢家
        address winner = address(0);
        if (card1 == card2) {
            // 平局，什么都不做
        } else if (card1 == CardType.Rock && card2 == CardType.Scissors ||
                   card1 == CardType.Paper && card2 == CardType.Rock ||
                   card1 == CardType.Scissors && card2 == CardType.Paper) {
            // 玩家1获胜
            winner = player1;
        } else {
            // 玩家2获胜
            winner = player2;
        }

        // 把两张牌放到赢家的牌后面，如果相同则丢弃
        if(winner != address(0)) {
            table.players[winner].cards.push(card1);
            table.players[winner].cards.push(card2);
        }
    
        // 比较当前局：判断游戏是否结束
        uint256 num1 = table.players[player1].cards.length;
        uint256 num2 = table.players[player2].cards.length;
        bool isFinish = true;
        address gameWinner = address(0);  // 当前局赢家
        if(0 == num1 && 0 == num2) {
            // 平局
        } else if(0 == num1) {
            // player2赢
            gameWinner = player2;
        } else if(0 == num2) {
            // player1赢
            gameWinner = player1;
        } else {
            isFinish = false;
        }

        if(isFinish) {
            emit GameOver(table.id, gameWinner);
            table.state = TableState.Over;
        }

         // 发送回合结果事件
        emit EndRound(table.id, round, winner, [card1, card2]);
    }
}
