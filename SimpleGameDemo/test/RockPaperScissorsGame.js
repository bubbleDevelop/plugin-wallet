const { expect } = require("chai");
const { network } = require("hardhat");

// 自由比较石头剪刀布

describe("RockPaperScissorsGame - Random Comparison", function () {
  // 同步标志：私链（单节点） -- 交易可以立即上链，可以立即查询交易回执；多节点 -- 交易上链会有延迟
  let syncFlag = true;
  // 使用异步方式的网络名称
  const asyncNetsFlag = ["hardhat", "ganache", "local"];
  let waitTxs = [];
  // 总账户私钥
  const privateKey = "0x02da90597bf4cef6621103622f27a31d65c0856a0a66ba2fd03e4663161f1c5b";
  let RockPaperScissorsGame;  // RockPaperScissorsGame合约对象
  let game;        // RockPaperScissorsGame部署对象
  let admin;
  let player1;
  let player2;

  // game
  let tableId = 0; // 记录tableId
  const playerNum = 2;
  let round = 0;  // 每局游戏当前的轮数

  const PreReady = 0;  // 0: 预准备状态
  const Ready = 1;     // 1：准备状态
  const Active = 2;    // 2: 游戏中
  const Over = 3;       // 3:游戏结束


  // 查看玩家信息
  async function getPlayerInfo(TableId, playerAddr) {
    console.log("getPlayerInfo=======");
    playerInfo = await game.getPlayerInfo(TableId, playerAddr);
    console.log("       player:",playerAddr, ", isInTable:", playerInfo.isInTable, ", cards:", playerInfo.cards);
    return playerInfo;
  }

  // 查看游戏桌信息
  async function getTableInfo(TableId) {
    console.log("getTableInfo, tableID:", TableId);
    table = await game.getTableInfo(TableId);
    console.log("       state:",table.state, ", round:", table.round, ", battleCards:", table.battleCards, ", playerAddrList:", table.playerAddrList);
    return table
  }

  // 等待交易
  async function waitTxList(EventName = "") {
    for (const tx of waitTxs) {
      if(EventName !== "") {
        console.log("EventName:", EventName, ", tx:", tx.hash);
      }
      await tx.wait();
    }
    // 清空交易列表
    waitTxs = [];
  }

  // 参与游戏
  async function joinTable(TableId, player) {
    // 游戏
    waitTxs.push(await game.connect(player).joinTable(TableId));
    if(syncFlag) {
      await waitTxList("joinTable");
    }
  }

  // 玩家2准备游戏
  async function readyGame(TableId, player) {
    // 游戏
    waitTxs.push(await game.connect(player).readyGame(TableId));
    if(syncFlag) {
      await waitTxList("readyGame");
    }
  }

  // 玩家1开始游戏
  async function startGame(TableId, player) {
    // 游戏
    waitTxs.push(await game.connect(player).startGame(TableId));
    if(syncFlag) {
      await waitTxList("startGame");
    }
  }

  // 出牌
  async function playCard(TableId, player, index) {
    // 游戏
    waitTxs.push(await game.connect(player).playCard(TableId, index));
    if(syncFlag) {
      await waitTxList("playCard");
    }
  }

  // 转账原生代币(同步模式)
  async function transferNative(toAddrList = []) {
    if(0 === toAddrList.length) {
      return;
    }
    console.log("transferNative===========================");
    if(!syncFlag) return;
    const provider = ethers.provider;
    console.log("provider============:", provider);
    var wallet = new ethers.Wallet(privateKey, provider);
    console.log("wallet============:", wallet);
    const totalBalance = await wallet.getBalance();
    console.log("totalBalance==========:", totalBalance);
    const transferAmount = ethers.utils.parseEther("1");
    if(parseInt(totalBalance, 10) === 0 || parseInt(totalBalance, 10) < transferAmount / toAddrList.length) {
      return;
    }
    // console.log("totalBalance:", parseInt(totalBalance, 10));
    // console.log("transferAmount:", transferAmount);
    for(const addr of toAddrList) {
      const balance = await provider.getBalance(addr);
      // console.log("balance:", balance);
      if(parseInt(balance, 10) === 0) 
      {
        let tx = await wallet.sendTransaction({
          // gasLimit: gasLimit,
          // gasPrice: gasPrice,
          to: addr,
          value: transferAmount
        });
        waitTxs.push(tx);
      }
    }

    await waitTxList("Transfer"); 
  }

  before(async function () {
    round = 0;
    for(i = 0; i < asyncNetsFlag.length; i++) {
      // 使用异步方法
      if(network.name === asyncNetsFlag[i]) {
        syncFlag = false;
        break;
      }
    }

    [admin, player1, player2] = await ethers.getSigners();
    // 部署RockPaperScissorsGame合约
    RockPaperScissorsGame = await ethers.getContractFactory("RockPaperScissorsGame");
    game = await RockPaperScissorsGame.deploy();
    
    // 转原生代币:player1, player2需要原生代币支付交易的手续费
    // transferNative([player1.address, player2.address]);
  });

  // 创建游戏桌
  describe("CreateTable", function () {
    it("Admin creates game table.", async function () {
      waitTxs.push(await game.createTable());
      if(syncFlag) {
        console.log("game address:", game.address);
        await waitTxList("CreateTable");
      }
    });
  });

  // 玩家参与游戏
  describe("JoinTable", function () {
    it("Player1 Join Table.", async function () {
      await joinTable(tableId, player1);
    });

    it("Player2 Join Table.", async function () {
      await joinTable(tableId, player2);
    });


    it("The number of players participating in the game should be 2", async function () {
      info = await game.getTableInfo(tableId);
      expect(info.playerAddrList.length).to.equal(playerNum);
    });
  });

  // 玩家2准备游戏
  describe("ReadyGame", function () {
    it("Player2 Ready Game.", async function () {
      await readyGame(tableId, player2);
    });
  });

  // 玩家1开始游戏
  describe("StartGame", function () {
    it("Player1 Start Game.", async function () {
      await startGame(tableId, player1);
    });
  });

  // 第一轮：玩家出牌
  describe("First round: Player PlayCard", function () {
    // player1出石头
    it("Player1 Play Card.", async function () {
      await playCard(tableId, player1, 0);
    });

    // player2出布
    it("Player2 Play Card.", async function () {
      await playCard(tableId, player2, 2);
    });
   

    // 比较玩家信息
    it("Compare Player1's Card and Player2's Card.", async function () {
      playerInfo = await getPlayerInfo(tableId, player1.address);
      expect(playerInfo.cards.length).to.equal(2);

      playerInfo = await getPlayerInfo(tableId, player2.address);
      expect(playerInfo.cards.length).to.equal(4);
    });
    
    // 获取Table信息
    it("Get the number of rounds in the table information.", async function () {
      tableInfo = await getTableInfo(tableId);
      expect(tableInfo.round).to.equal(1);
    });
  });

  // 第二轮：玩家出牌
  describe("Second round: Player PlayCard", function () {
    // player1出剪刀
    it("Player1 Play Card.", async function () {
      await playCard(tableId, player1, 0);
    });

    // player2出石头
    it("Player2 Play Card.", async function () {
      await playCard(tableId, player2, 0);
    });

    // 比较玩家信息
    it("Compare Player1's Card and Player2's Card.", async function () {
      playerInfo = await getPlayerInfo(tableId, player1.address);
      expect(playerInfo.cards.length).to.equal(1);

      playerInfo = await getPlayerInfo(tableId, player2.address);
      expect(playerInfo.cards.length).to.equal(5);
    });
    
    // 获取Table信息
    it("Get the number of rounds in the table information.", async function () {
      tableInfo = await getTableInfo(tableId);
      expect(tableInfo.round).to.equal(2);
    });
  });

  // 第三轮：玩家出牌
  describe("Third round: Player PlayCard, End Round", function () {
    // player1出布
    it("Player1 Play Card.", async function () {
      await playCard(tableId, player1, 0);
    });

    // player2出剪刀
    it("Player2 Play Card, End Round, Player2 Winner.", async function () {
      // 解析交易回执中的事件
      const provider = ethers.provider;
      txReceipt = await game.connect(player2).playCard(tableId, 0);
      if(syncFlag) {
        waitTxs.push(txReceipt);
        await waitTxList("playCard and EndRound");
      }
      // 获取交易回执
      const receipt = await provider.getTransactionReceipt(txReceipt.hash);
      let gameOver = false;
      // 遍历事件日志
      for (const log of receipt.logs) {
        // 判断事件名称
        const logInfo = game.interface.parseLog(log);
        // console.log("logInfo:", logInfo);
        let logMsg = "";
        if (logInfo.name === "EndRound") {
          // 当前轮

        } else if(logInfo.name === "GameOver") {
          // console.log(logInfo);
          gameOver = true;
          // winner is player2
          expect(logInfo.args._winner).to.equal(player2.address);
        }
      }
      
      // 游戏结束
      expect(gameOver).to.equal(true);

    });

    // 比较玩家信息
    it("Compare Player1's Card and Player2's Card.", async function () {
      playerInfo = await getPlayerInfo(tableId, player1.address);
      expect(playerInfo.cards.length).to.equal(0);

      playerInfo = await getPlayerInfo(tableId, player2.address);
      expect(playerInfo.cards.length).to.equal(6);
    });
    
    // 获取Table信息
    it("Get the number of rounds in the table information.", async function () {
      tableInfo = await getTableInfo(tableId);
      expect(tableInfo.round).to.equal(3);
    });
  });

  // 游戏结束
  describe("Game Over", function () {
    // 游戏结束，验证清结算结果
    it("The game is over", async function () {
      // 清结算
      tableInfo = await game.getTableInfo(tableId);
      expect(tableInfo.state).to.equal(Over);
    });
  });

  // 代签代扣
    describe("setIssuer, setLineOfCredit be success", function (accounts) {
      it("check setIssuer and setLineOfCredit", async function () {
        const { ethers } = require("hardhat");
        const [owner, addr1, addr2] = await ethers.getSigners();
        const game = await ethers.deployContract("RockPaperScissorsGame");

        var issuerAddres = await game.issuer.call();
        expect(issuerAddres).to.equal(owner.address);

        await game.setIssuer(addr1.address);
        issuerAddres = await game.issuer.call();
        expect(issuerAddres).to.equal(addr1.address);

        const amount = 123456;
        await game.connect(addr1).setLineOfCredit(amount)
        const lineOfCredit = await game.lineOfCredit.call();
        expect(amount).to.equal(lineOfCredit.toNumber());
      });
    });
});
