let game
var Game = artifacts.require("Game");

contract('Game', (accounts) => {
    beforeEach(async () => {
        game = await Game.new();
        console.log('game address:', game.address);
    })

    it('Game: setIssuer and setLineOfCredit be success', async () => {
        await game.setIssuer(accounts[1]);
        var issuerAddres = await game.issuer.call();
        console.log('issuer address:',issuerAddres);
        assert.strictEqual(issuerAddres, accounts[1]);

        const amount = 123456;
        await game.setLineOfCredit(amount)
        const lineOfCredit = await game.lineOfCredit.call();
        console.log('lineOfCredit value:', lineOfCredit.toNumber());
        assert.strictEqual(amount, lineOfCredit.toNumber())
    });

});