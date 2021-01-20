const express = require('express');
const { update } = require('../../models/Item');
const router = express.Router();

// Item Model
const Item = require('../../models/Item');


async function findSKU (req, res) {
    let sku
    try {
        sku = await Item.findById(req.params.id)
        if(sku == null){
            return res.status(404).json({ message : 'Cannot find SKU'})
        }
    } catch (error) {
        console.log(error)
        res.json(error)
    }
    res.sku = sku

}

// @route GET api.items
// @desc Get All Items
// @access Public
router.get('/', (req, res) => {
 //   console.log('Items Retrieved')
    Item.find()
    .sort({date : -1 })
    .then(items => res.json(items));
});

// @route POST api.items
// @desc Create a Item
// @access Public
router.post('/', (req, res) => {
    const newItem = new Item({
        name: req.body.name,
        status: req.body.status
    });

    newItem.save().then(item => res.json(item)).catch(err => {
        res.json({message : err})
    }) ;
    console.log('New Item Created')
});

// @route DELETE api.items/:id
// @desc Delete An Item
// @access Public
router.delete('/:id', (req, res) => {
Item.findById(req.params.id)
.then(item => item.remove().then(() => res.json({success: true, message : ({'Removed Item' : req.params.id})})))
.catch(err => res.status(404).json({success: false}));
});

// @route UPDATE api.items/:id
// @desc Update An Item
// @access Public

router.patch('/:id', async (req, res) => {
    let actualId = req.params.id
    console.log(actualId)
    console.log('Someone Wants to Update')
    console.log(req.body)
    try {
        const updatedSku = await Item.findByIdAndUpdate(actualId, req.body);
        console.log(updatedSku)
        res.json(updatedSku)
    } catch (error) {
        res.status(400).json({ message : error.message })
        console.log(error)
    }
   
})
// Find


module.exports = router