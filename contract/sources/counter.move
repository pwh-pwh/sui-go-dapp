
module counter::counter;

public struct Counter has key {
    id: UID,
    counter: u8
}



fun init(ctx : &mut TxContext) {

    let counter = Counter {
        id: object::new(ctx),
        counter: 0
    };
    transfer::share_object(counter);
}


entry fun increment(counter:&mut Counter,_ctx: &mut TxContext) {
    counter.counter = counter.counter + 1;
}

entry fun addNum(counter:&mut Counter, num: u8, _ctx: &mut TxContext) {
    counter.counter = counter.counter + num;
}

// package id:0xbda58f110ce755a63c007d68cc53f7ac68c780dc8fb1fb16ad52d797143b4799
//tx: 6S5b62crgfbikAEZKLwphkQqtf5wuJwQqkkRW4ywnxug