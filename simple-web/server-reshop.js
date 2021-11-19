// ExpressJS Setup
const express = require('express');
const app = express();
var bodyParser = require('body-parser');
const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');

// Hyperledger Bridge connection.json가져와서 구조체화
const ccpPath = path.resolve(__dirname,'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

// Constants
const PORT = 3000;
const HOST = '0.0.0.0';

// use static file
app.use(express.static(path.join(__dirname, 'views')));

// configure app to use body-parser
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

// main page routing
app.get('/', (req, res)=>{
    res.sendFile(__dirname + '/index-reshop.html');
});

// repair issue // 생성
app.post('/repair', async(req, res)=>{
    const mode = req.body.mode;
// register, respond, request, complete, pay
    if(mode == 'register')
    {
        const coid = req.body.coid;
        const cuid = req.body.cuid;
        const carinfo = req.body.carinfo;
        result = cc_call('register', [coid, cuid, carinfo], res);
    }else if(mode == 'respond'){
        const coid = req.body.coid;
        const sid = req.body.sid;
        const items = req.body.items;
        const price = req.body.price;
        result = cc_call('respond', [coid, sid, items, price], res);
    }else if(mode == 'request'){
        const coid = req.body.coid;
        const cuid = req.body.cuid;
        result = cc_call('request', [coid, cuid], res);
    }else if(mode == 'complete'){
        const coid = req.body.coid;
        const sid = req.body.sid;
        result = cc_call('complete', [coid, sid], res);
    }else if(mode == 'pay'){
        const coid = req.body.coid;
        const cuid = req.body.cuid;
        result = cc_call('pay', [coid, cuid], res);
    }
});

// history
app.get('/repair', async(req, res)=>{

    try {
        const id = req.query.coid;
        console.log(`Contract history query: ${id}`);

        cc_call('history', [id], res)
    }
    catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        //process.exit(1);
    }
});

async function cc_call(fn_name, args, respond){

    // 지갑경로확인 유무확인
    const walletPath = path.join(process.cwd(), 'wallet');
    const wallet = new FileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    const userExists = await wallet.exists('user1');
    if (!userExists) {
        console.log(`cc_call`);
        console.log('An identity for the user "user1" does not exist in the wallet');
        console.log('Run the registerUser.js application before retrying');
        return;
    }
    // 게이트웨이 연결
    const gateway = new Gateway();
    await gateway.connect(ccpPath, { wallet, identity: 'user1', discovery: { enabled: true, asLocalhost: true } });
    const network = await gateway.getNetwork('rechannel'); // 채널연결
    const contract = network.getContract('reshop'); // 체인코드 연결

    var result;
    // register, respond, request, complete, pay, history 
    if(fn_name == 'register'){
        result = await contract.submitTransaction('register', args[0],args[1],args[2]);
    }else if(fn_name == 'respond'){
        result = await contract.submitTransaction('respond', args[0],args[1],args[2],args[3]);
    }else if(fn_name == 'request'){
        result = await contract.submitTransaction('request', args[0],args[1]);
    }else if(fn_name == 'complete'){
        result = await contract.submitTransaction('complete', args[0],args[1]);
    }else if(fn_name == 'pay'){
        result = await contract.submitTransaction('pay', args[0],args[1]);
    }else if(fn_name == 'history'){
        result = await contract.evaluateTransaction('history', args[0]);
        var myobj = JSON.parse(result);
        respond.status(200).json(myobj);
        return;
        //console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
    }else{
        result = 'not supported function'
    }

    gateway.disconnect();

    result = '{"result":"tx has been submitted"}';
    var stobj = JSON.parse(result);
    respond.status(200).json(stobj)
    return;
}

// server start
app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);