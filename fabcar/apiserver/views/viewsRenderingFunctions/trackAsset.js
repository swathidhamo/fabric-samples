const { FileSystemWallet, Gateway } = require('fabric-network');
const path = require('path');
const ccpPath = path.resolve(__dirname,'..','..',  '..', '..', 'first-network', 'connection-org1.json');


async function trackAsset(req,res){
  try{
  	
  	const walletPath = path.join(process.cwd(), 'wallet');
  	const wallet = new FileSystemWallet(walletPath);
  	console.log(`Wallet path: ${walletPath}`);

  	const userExists = await wallet.exists('user1');
  	if (!userExists) {
  		console.log('An identity for the user "user1" does not exist in the wallet');
  		console.log('Run the registerUser.js application before retrying');
  		return;	
  	 }
  	const gateway = await new Gateway();		
  	await gateway.connect(ccpPath, { wallet, identity: 'user1', discovery: { enabled: true, asLocalhost: true } });

  	const network = await gateway.getNetwork('mychannel');
  	const contract = network.getContract('fabcar');
  	const assetID = "ASSET_" + req.params.assetid;
  	const result = await contract.evaluateTransaction('getReadingForID', assetID);
  	console.log(`Transaction has been evaluated, result is: ${result}`);
  	var resultString = result.toString();
  	var jsonData = JSON.parse(resultString);
  	console.log(jsonData.toString() + "Latitute");
  	res.render('trackAsset',{
  			"assetid": req.params.assetid,
  			"assetHistory":jsonData
  	});


  }
  catch(error) {
  	console.error(`Failed to evaluate transaction: ${error}`);
  	res.status(500).json({error: error});
  	process.exit(1);
  }
   
}

module.exports = trackAsset;
