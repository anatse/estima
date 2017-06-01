function foo(params) {
    let db = require('internal').db,
        toCol = db._collection (params.toColName),
        edgeCol = db._collection(params.edgeColName),
        doc = db._collection(params.fromColName).document(params.fKey);

    // Save new comment
    let toDoc = toCol.save (params.comment);

    // Add edge between comment and object
    edgeCol.save ({_from: doc._id, _to: toDoc._id, label: 'comment'});

    // Add edge between comment and user
    edgeCol.save ({_from: params.userId, _to: toDoc._id, label: 'userComment'});

    return {success: true, entityKey: toDoc._key};
}
