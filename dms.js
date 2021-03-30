const {
	targetMonitor
} = require('./sites/targetEfficent');
const {
	newEggMonitor
} = require('./sites/newEgg')
const {
	gameStopMonitor
} = require('./sites/gameStop');
const { bestBuyMonitor } = require('./sites/bestBuy')

const { amdMonitor } = require('./sites/amd')
const { amdSiteMonitor } = require('./sites/amdSite')
const { walmartMonitor } = require('./sites/walmart')
const delay = require('delay');
require('newrelic');
var skuBank = []

const rp = require('request-promise').defaults({
	followAllRedirects: true,
	resolveWithFullResponse: true,
	gzip: true,
});

function SKUADD(clients, triggerText, replyText) {
	try {
		clients.on('message', async (message) => {
			if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
				// message.author.send(replyText);
				let pricerange = ''
				const content = message.content;
				const site = content.split(' ')[1];
				const SKU = content.split(' ')[2];
				if(content.includes('[')){
					pricerange = content.split('[')[1].split(']')[0]

				}
				//    fetch('')
				console.log(site)
				console.log(SKU)
				console.log(content)
				console.log(pricerange)
				
				if (SKU.length > 1 && site.length > 1) {
					if (site.toUpperCase() == 'TARGET') {
						skuBank.push({
							sku: SKU,
							site: 'TARGET',
							stop: false
						})
						let monitor = new targetMonitor(SKU.toString())
						monitor.task()
					} else if (site.toUpperCase() == 'NEWEGG') {
						skuBank.push({
							sku: SKU,
							site: 'NEWEGG',
							stop: false
						})
						let monitor = new newEggMonitor(SKU.toString())
						monitor.task()
					} else if (site.toUpperCase() == 'GAMESTOP') {
						skuBank.push({
							sku: SKU,
							site: 'GAMESTOP',
							stop: false
						})
						let monitor = new gameStopMonitor(SKU.toString())
						monitor.task()
					} else if (site.toUpperCase() == 'AMD') {
						skuBank.push({
							sku: SKU,
							site: 'AMD',
							stop: false
						})
						let monitor = new amdMonitor(SKU.toString())
						monitor.task()
					} else if (site.toUpperCase() == 'AMDSITE') {
						skuBank.push({
							sku: SKU,
							site: 'AMDSITE',
							stop: false
						})
						let monitor = new amdSiteMonitor(SKU.toString())
						monitor.task()
					} else if (site.toUpperCase() == 'WALMART') {
						skuBank.push({
							sku: SKU,
							site: 'WALMART',
							stop: false
						})
						let monitor = new walmartMonitor(SKU.toString(), pricerange)
						monitor.task()
					} else if (site.toUpperCase() == 'BESTBUY') {
						skuBank.push({
							sku: SKU,
							site: 'BESTBUY',
							stop: false
						})
						let monitor = new bestBuyMonitor(SKU.toString())
						monitor.task()
					}
				console.log(skuBank)
				message.channel.send(`${SKU} Added to ${site}`)
			}
		}});
	} catch (error) {
		console.log(error);
	}

}

function findCommand(clients, triggerText, replyText) {
	clients.on('message', message => {
		if (message.content.toLowerCase() === triggerText.toLowerCase()) {
			message.author.send(replyText);
		}
	});

}

function deleteSku(clients, triggerText, replyText) {
	try {
		clients.on('message', async (message) => {
			if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
				// message.author.send(replyText);
				const content = message.content;
				const site = content.split(' ')[1];
				const SKU = content.split(' ')[2];
				console.log(site)
				console.log(`SKU - ${SKU}`)
				console.log(content)
				let index = skuBank.findIndex(e => e.sku == SKU)
				// console.log(index)
				// console.log(skuBank[index])
				skuBank[index].stop = true;
				// console.log(skuBank)
				(async () => {
					await delay(10000)
					skuBank.splice(index, 1)
				})()
				
				// console.log(skuBank)
				function replaceWithTheCapitalLetter(values){
				return values.charAt(0).toUpperCase() + values.slice(1);
				}
				message.channel.send(`${SKU} Deleted From ${replaceWithTheCapitalLetter(site)}`)
				return;
				//    fetch('')

			}
		});
	} catch (error) {
		console.log(error);
	}
}

function checkBank (clients, triggerText, replyText){
	try {
		clients.on('message', async (message) => {
			if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
				let string = (skuBank)
			//	message.channel.send(JSON.parse(string))
				for (let i = 0; i < skuBank.length; i++){
					console.log(skuBank[i])
					message.channel.send(JSON.stringify(skuBank[i]))
					await delay(1000)
				}
			//	message.channel.send(string)
			}
		});
	} catch (error) {
		console.log(error);
	}
}

function massAdd (clients, triggerText, replyText){
	try {
		clients.on('message', async (message) => {
			if (message.content.toLowerCase().includes(triggerText.toLowerCase())) {
				let string = message.content
				const content = message.content;
				const site = content.split(' ')[1].split('|')[0]
				console.log(site.toUpperCase())
			//	const SKU = content.split(' ')[2];
			//	console.log(site)
				let g  = string.split('\n')
			//	console.log(g)
			for(let i = 0; i < g.length; i++){
				if(!g[i].includes('!massAdd')){
					console.log(g[i])
					console.log(site.toUpperCase())
					if (site.toUpperCase() == 'TARGET') {
						skuBank.push({
							sku: g[i],
							site: 'TARGET',
							stop: false
						})
						let monitor = new targetMonitor(SKU.toString())
						monitor.task()
					} else if (site.toUpperCase() == 'NEWEGG') {
						skuBank.push({
							sku: g[i],
							site: 'NEWEGG',
							stop: false
						})
						let monitor = new newEggMonitor(SKU.toString())
						monitor.task()
					} else if (site.toUpperCase() == 'GAMESTOP') {
						skuBank.push({
							sku: g[i],
							site: 'GAMESTOP',
							stop: false
						})
						let monitor = new gameStopMonitor(SKU.toString())
						monitor.task()
					} else if (site.toUpperCase() == 'AMD') {
						skuBank.push({
							sku: g[i],
							site: 'AMD',
							stop: false
						})
						let monitor = new amdMonitor(SKU.toString())
						monitor.task()
					} else if (site.toUpperCase() == 'AMDSITE') {
						skuBank.push({
							sku: g[i],
							site: 'AMDSITE',
							stop: false
						})
						let monitor = new amdSiteMonitor(SKU.toString())
						monitor.task()
					} else if (site.toUpperCase() == 'WALMART') {
						skuBank.push({
							sku: g[i],
							site: 'WALMART',
							stop: false
						})
						let monitor = new walmartMonitor(SKU.toString())
						monitor.task()
					} else if (site.toUpperCase() == 'BESTBUY') {
						skuBank.push({
							sku: g[i],
							site: 'BESTBUY',
							stop: false
						})
						let monitor = new bestBuyMonitor(SKU.toString())
						
					}

				}
				await delay(30000)
			}
			message.channel.send("SKUS Added")
			}
		});
	} catch (error) {
		console.log(error);
	}
}


module.exports = {
	SKUADD,
	findCommand,
	deleteSku,
	checkBank,
	skuBank,
	massAdd
}