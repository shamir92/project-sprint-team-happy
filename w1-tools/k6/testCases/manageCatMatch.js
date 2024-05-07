import { fail } from 'k6';
import {
  generateRandomDescription, generateRandomImageUrl, generateRandomNumber, testDelete, testGet, testPutJson,

  generateTestObjects, generateUniqueName, isExists, testPostJson, assert,
} from '../helper.js';
import { generateRandomCatBreed, generateRandomCatGender } from '../types/cat.js';

const TEST_NAME = '(manage cat match)';
/**
 * Test manage cat match POST API
 * @param {Config} config
 * @param {Object} tags
 * @param {User} user
 */
export function TestPostManageCatMatch(config, user, user2, tags = {}) {
  let res, currentTest;
  // eslint-disable-next-line no-undef
  const route = `${__ENV.BASE_URL}/v1/cat/match`;
  const currentFeature = `${TEST_NAME} | post manage cat match`;
  if (!user) fail(`${currentFeature} fail due to user is empty`);
  const headers = {
    Authorization: `Bearer ${user.accessToken}`,
  };

  const [ownedCats, notOwnedCats, payload] = generateCatMatch(config, currentFeature, headers, {
    Authorization: `Bearer ${user2.accessToken}`,
  }, tags);

  const manageCatMatchNegativePayloads = generateTestObjects({
    matchCatId: { type: 'string', notNull: true },
    userCatId: { type: 'string', notNull: true },
    message: {
      type: 'string', notNull: true, minLength: 1, maxLength: 120,
    },
  }, {
    matchCatId: notOwnedCats[generateRandomNumber(0, notOwnedCats.length - 1)].id,
    userCatId: ownedCats[generateRandomNumber(0, ownedCats.length - 1)].id,
    message: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
  });

  if (!config.POSITIVE_CASE) {
    currentTest = 'no header';
    res = testPostJson(route, {}, {}, tags, ['noContentType']);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 401`]: (r) => r.status === 401,
    });


    currentTest = 'no body';
    res = testPostJson(route, {}, headers, tags, ['noContentType']);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 400`]: (r) => r.status === 400,
    });


    currentTest = 'invalid payload';
    manageCatMatchNegativePayloads.forEach((payload) => {
      res = testPostJson(route, payload, headers, tags);
      assert(res, currentFeature, config, {
        [`${currentTest} should return 400`]: (r) => r.status === 400,
      }, payload);
    });


    currentTest = 'match a cat that the id is not exist';
    res = testPostJson(route, {
      matchCatId: '123456789012345678901234',
      userCatId: ownedCats[0].id,
      message: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
    }, headers, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 404`]: (r) => r.status === 404,
    }, {
      matchCatId: '123456789012345678901234',
      userCatId: ownedCats[0].id,
      message: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
    });


    currentTest = 'match cat that is both owned by user';
    res = testPostJson(route, {
      matchCatId: ownedCats[0].id,
      userCatId: ownedCats[1].id,
      message: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
    }, headers, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 400`]: (r) => r.status === 400,
    }, {
      matchCatId: ownedCats[0].id,
      userCatId: ownedCats[1].id,
      message: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
    });


    currentTest = 'match a cat that is not owned';
    res = testPostJson(route, {
      matchCatId: notOwnedCats[0].id,
      userCatId: notOwnedCats[1].id,
      message: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
    }, headers, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 404`]: (r) => r.status === 404,
    }, {
      matchCatId: notOwnedCats[0].id,
      userCatId: notOwnedCats[1].id,
      message: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
    });


    currentTest = 'match a cat that is the same gender';
    const randomGender = generateRandomCatGender();
    res = testPostJson(route, {
      matchCatId: notOwnedCats.find((cat) => cat.sex === randomGender).id,
      userCatId: ownedCats.find((cat) => cat.sex === randomGender).id,
      message: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
    }, headers, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 400`]: (r) => r.status === 400,
    }, {
      matchCatId: notOwnedCats.find((cat) => cat.sex === randomGender).id,
      userCatId: ownedCats.find((cat) => cat.sex === randomGender).id,
      message: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
    });
  }

  if (!config.POSITIVE_CASE) {
    currentTest = 'match a cat that already matched';
    const exRes = testPostJson(route, payload, headers, tags);
    assert(exRes, currentFeature, config, {
      [`${currentTest} should return 400`]: (r) => r.status === 400,
    }, payload);
  }
}

/**
 * Test manage cat match GET API
 * @param {Config} config
 * @param {Object} tags
 * @param {User} user
 * @param {import('../types/cat.js').Cat} cat
 * @returns {CatMatch | null}
 */
export function TestGetManageCatMatch(config, user, tags = {}) {
  let res, currentTest;
  // eslint-disable-next-line no-undef
  const route = `${__ENV.BASE_URL}/v1/cat/match`;
  // eslint-disable-next-line no-undef
  const getRoute = `${__ENV.BASE_URL}/v1/cat`;
  const currentFeature = `${TEST_NAME} | get manage cat match`;
  if (!user) fail(`${currentFeature} fail due to user is empty`);

  const headers = {
    Authorization: `Bearer ${user.accessToken}`,
  };

  if (!config.POSITIVE_CASE) {
    currentTest = 'no header';
    res = testGet(route, {}, {}, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 401`]: (r) => r.status === 401,
    });
  }

  const commonChecks = (currFeat) => ({
    [`${currFeat} should return 200`]: (r) => r.status === 200,
    [`${currFeat} should have all id`]: (r) => isExists(r, 'data.[].id'),
    // [`${currFeat} should have all issuedBy.name`]: (r) => isExists(r, 'data.[].issuedBy.name'),
    // [`${currFeat} should have all issuedBy.email`]: (r) => isExists(r, 'data.[].issuedBy.email'),
    // [`${currFeat} should have all issuedBy.createdAt`]: (r) => isExists(r, 'data.[].issuedBy.createdAt'),
    // [`${currFeat} should have all matchCatDetail.id`]: (r) => isExists(r, 'data.[].matchCatDetail.id'),
    // [`${currFeat} should have all matchCatDetail.name`]: (r) => isExists(r, 'data.[].matchCatDetail.name'),
    // [`${currFeat} should have all matchCatDetail.race`]: (r) => isExists(r, 'data.[].matchCatDetail.race'),
    // [`${currFeat} should have all matchCatDetail.sex`]: (r) => isExists(r, 'data.[].matchCatDetail.sex'),
    // [`${currFeat} should have all matchCatDetail.description`]: (r) => isExists(r, 'data.[].matchCatDetail.description'),
    // [`${currFeat} should have all matchCatDetail.ageInMonth`]: (r) => isExists(r, 'data.[].matchCatDetail.ageInMonth'),
    // [`${currFeat} should have all matchCatDetail.imageUrls`]: (r) => isExists(r, 'data.[].matchCatDetail.imageUrls'),
    // [`${currFeat} should have all matchCatDetail.hasMatched`]: (r) => isExists(r, 'data.[].matchCatDetail.hasMatched'),
    // [`${currFeat} should have all matchCatDetail.createdAt`]: (r) => isExists(r, 'data.[].matchCatDetail.createdAt'),
    // [`${currFeat} should have all userCatDetail.id`]: (r) => isExists(r, 'data.[].userCatDetail.id'),
    // [`${currFeat} should have all userCatDetail.name`]: (r) => isExists(r, 'data.[].userCatDetail.name'),
    // [`${currFeat} should have all userCatDetail.race`]: (r) => isExists(r, 'data.[].userCatDetail.race'),
    // [`${currFeat} should have all userCatDetail.sex`]: (r) => isExists(r, 'data.[].userCatDetail.sex'),
    // [`${currFeat} should have all userCatDetail.description`]: (r) => isExists(r, 'data.[].userCatDetail.description'),
    // [`${currFeat} should have all userCatDetail.ageInMonth`]: (r) => isExists(r, 'data.[].userCatDetail.ageInMonth'),
    // [`${currFeat} should have all userCatDetail.imageUrls`]: (r) => isExists(r, 'data.[].userCatDetail.imageUrls'),
    // [`${currFeat} should have all userCatDetail.hasMatched`]: (r) => isExists(r, 'data.[].userCatDetail.hasMatched'),
    // [`${currFeat} should have all userCatDetail.createdAt`]: (r) => isExists(r, 'data.[].userCatDetail.createdAt'),
    [`${currFeat} should have all message`]: (r) => isExists(r, 'data.[].message'),
    [`${currFeat} should have all createdAt`]: (r) => isExists(r, 'data.[].createdAt'),
  });
  currentTest = 'get all match cats';
  res = testGet(route, {}, headers, tags);
  assert(res, currentFeature, config, commonChecks(currentTest));

  currentTest = 'get all cats that owned and not matched';
  res = testGet(getRoute, { owned: true, hasMatched: false, limit: 1000, offset: 0 }, headers, tags);
  assert(res, currentFeature, config, {
    [`${currentTest} should return 200`]: (r) => r.status === 200,
  }, { owned: true });
  /** @type {Cat[]} */
  const ownedCats = res.json().data;
  const userCat = ownedCats[generateRandomNumber(0, ownedCats.length - 1)]


  currentTest = 'get all cats that is not owned, not matched and have opposite gender';
  res = testGet(getRoute, { owned: false, hasMatched: false, sex: generateRandomCatGender(userCat.sex), limit: 1000, offset: 0 }, headers, tags);
  assert(res, currentFeature, config, {
    [`${currentTest} should return 200`]: (r) => r.status === 200,
  }, { owned: false, limit: 1000, offset: 0 });
  /** @type {Cat[]} */
  const notOwnedCats = res.json().data;
  const matchCat = notOwnedCats[generateRandomNumber(0, notOwnedCats.length - 1)]

  currentTest = 'match a new cat';
  const positivePayload = {
    matchCatId: matchCat.id,
    userCatId: userCat.id,
    message: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
  };
  res = testPostJson(route, positivePayload, headers, tags);
  assert(res, currentFeature, config, {
    [`${currentTest} should return 201`]: (r) => r.status === 201,
  }, positivePayload);

  if (!config.POSITIVE_CASE) {
    currentTest = 'cat that is matched should not be able to edit the gender'
    // eslint-disable-next-line no-undef
    res = testPutJson(`${__ENV.BASE_URL}/v1/cat/${userCat.id}`, {
      name: generateUniqueName(),
      race: generateRandomCatBreed(),
      ageInMonth: generateRandomNumber(1, 120082),
      sex: generateRandomCatGender(userCat.sex),
      description: generateRandomDescription(200),
      imageUrls: [generateRandomImageUrl()],
    }, headers, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 400`]: (r) => r.status === 400,
    });
  }
}

/**
 * Test manage cat match DELETE API
 * @param {Config} config
 * @param {User} user
 * @param {Object} tags
 */
export function TestDeleteManageCatMatch(config, user, user2, tags = {}) {
  let res, currentTest;

  const headers = {
    Authorization: `Bearer ${user.accessToken}`,
  };
  // eslint-disable-next-line no-undef
  const getRoute = `${__ENV.BASE_URL}/v1/cat/match`;
  const currentFeature = `${TEST_NAME} | delete manage cat`;
  if (!user) fail(`${currentFeature} fail due to user is empty`);

  generateCatMatch(config, currentFeature, headers, {
    Authorization: `Bearer ${user2.accessToken}`,
  }, tags, true);

  currentTest = 'get all match cats';
  res = testGet(getRoute, {}, headers, tags);
  assert(res, currentFeature, config, {
    [`${currentTest} should return 200`]: (r) => r.status === 200,
  });
  let catMatch = res.json().data.find(
    /** @param {CatMatch} match */
    (match) => match.userCatDetail.hasMatched === false
      && match.issuedBy.email == user.email
      && match.matchCatDetail.hasMatched === false);


  // eslint-disable-next-line no-undef
  const route = `${__ENV.BASE_URL}/v1/cat/match/${catMatch.id}`;

  if (!config.POSITIVE_CASE) {
    // Negative case, no header
    currentTest = 'no header';
    res = testDelete(route, {}, {}, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 401`]: (r) => r.status === 401,
    });
  }

  // Positive case, delete cat
  currentTest = 'delete match cat'
  res = testDelete(route, {}, headers, tags);
  let positivePayloadPassAssertTest = assert(res, currentFeature, config, {
    [`${currentTest} should return 200`]: (r) => r.status === 200,
  });

  // Positive case, deleted cat should not appear in the list
  currentTest = 'get all match cats';
  res = testGet(getRoute, {}, headers, tags);
  positivePayloadPassAssertTest = assert(res, currentFeature, config, {
    [`'${currentTest} should return 200`]: (r) => r.status === 200,
    [`'${currentTest} should not have deleted match request`]: (r) => {
      try {
        return r.json().data.find((match) => match.id === catMatch.id) === undefined;
      } catch (e) {
        return false;
      }
    },
  }, {});

  if (!positivePayloadPassAssertTest) return null;

  return res.json().data;
}

/**
 * Generate new cat match
 * @param {Config} config 
 * @param {string} currentFeature 
 * @param {Object} otherUserHeader 
 * @param {Object} tags 
 */
function generateCatMatch(config, currentFeature, userHeader, otherUserHeader, tags, generateByUser = false) {
  // eslint-disable-next-line no-undef
  const route = `${__ENV.BASE_URL}/v1/cat/match`;
  // eslint-disable-next-line no-undef
  const getRoute = `${__ENV.BASE_URL}/v1/cat`;


  const matchCatGender = generateRandomCatGender()
  const userCatGender = generateRandomCatGender(matchCatGender)

  let currentTest = 'get all cats not owned';
  let res = testGet(getRoute, { owned: true, limit: 1000, offset: 0 }, otherUserHeader, tags);
  assert(res, currentFeature, config, {
    [`${currentTest} should return 200`]: (r) => r.status === 200,
  }, { owned: true, otherUserHeader });
  /** @type {Cat[]} */
  let notOwnedCats = res.json().data;
  if (notOwnedCats.length === 0) {
    for (let i = 0; i < 2; i++) {
      const positivePayload = {
        name: generateUniqueName(),
        race: generateRandomCatBreed(),
        ageInMonth: generateRandomNumber(1, 120082),
        sex: matchCatGender,
        description: generateRandomDescription(200),
        imageUrls: [generateRandomImageUrl()],
      };
      currentTest = 'create new cat';
      res = testPostJson(route, positivePayload, otherUserHeader, tags);
      assert(res, currentFeature, config, {
        [`${currentTest} should return 201`]: (r) => r.status === 201,
      }, Object.assign(positivePayload, { otherUserHeader }));
    }
    let currentTest = 'get all cats not owned';
    let res = testGet(getRoute, { owned: true, limit: 1000, offset: 0 }, otherUserHeader, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 200`]: (r) => r.status === 200,
    }, { owned: true, otherUserHeader });
    /** @type {Cat[]} */
    notOwnedCats = res.json().data;
  }
  notOwnedCats = notOwnedCats.filter((cat) => cat.hasMatched === false)


  currentTest = 'get all cats that is owned';
  res = testGet(getRoute, { owned: true, limit: 1000, offset: 0 }, userHeader, tags); // todo: kebalik ini harusnya yg owned
  assert(res, currentFeature, config, {
    [`${currentTest} should return 200`]: (r) => r.status === 200,
  }, { owned: true, limit: 1000, offset: 0, userHeader });
  /** @type {Cat[]} */
  let ownedCats = res.json().data;
  ownedCats = ownedCats.filter((cat) => cat.hasMatched === false)
  let positivePayload
  if (!generateByUser) {
    currentTest = 'match a new cat';
    const notHasMatchedNotOwnedCat = notOwnedCats.filter((cat) => cat.hasMatched === false && cat.sex == matchCatGender)
    const notHasMatchedOwnedCat = ownedCats.filter((cat) => cat.hasMatched === false && cat.sex == userCatGender)
    positivePayload = {
      matchCatId: notHasMatchedOwnedCat[generateRandomNumber(0, notHasMatchedOwnedCat.length - 1)].id,
      userCatId: notHasMatchedNotOwnedCat[generateRandomNumber(0, notHasMatchedNotOwnedCat.length - 1)].id,
      message: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
    };
    res = testPostJson(route, positivePayload, otherUserHeader, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 201`]: (r) => r.status === 201,
    }, positivePayload);
  } else {
    currentTest = 'match a new cat';
    const notHasMatchedNotOwnedCat = notOwnedCats.filter((cat) => cat.hasMatched === false && cat.sex == matchCatGender)
    const notHasMatchedOwnedCat = ownedCats.filter((cat) => cat.hasMatched === false && cat.sex == userCatGender)
    positivePayload = {
      matchCatId: notHasMatchedNotOwnedCat[generateRandomNumber(0, notHasMatchedNotOwnedCat.length - 1)].id,
      userCatId: notHasMatchedOwnedCat[generateRandomNumber(0, notHasMatchedOwnedCat.length - 1)].id,
      message: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
    };
    res = testPostJson(route, positivePayload, userHeader, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 201`]: (r) => r.status === 201,
    }, positivePayload);

  }

  return [ownedCats, notOwnedCats, positivePayload]
}


/**
 * Test manage cat match approve POST API
 * @param {Config} config
 * @param {User} user
 * @param {Object} tags
 */
export function TestPostManageCatApprove(config, user, user2, tags = {}) {
  let res, currentTest;
  // eslint-disable-next-line no-undef
  const getRoute = `${__ENV.BASE_URL}/v1/cat/match`;
  // eslint-disable-next-line no-undef
  const route = `${__ENV.BASE_URL}/v1/cat/match/approve`;
  const currentFeature = `${TEST_NAME} | manage cat match approve`;
  if (!user) fail(`${currentFeature} fail due to user is empty`);

  const headers = {
    Authorization: `Bearer ${user.accessToken}`,
  };

  if (!config.POSITIVE_CASE) {
    // Negative case, no header
    currentTest = 'no header';
    res = testPostJson(route, {}, {}, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 401`]: (r) => r.status === 401,
    });

    // Negative case, no body
    currentTest = 'no body';
    res = testPostJson(route, {}, headers, tags, ['noContentType']);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 400`]: (r) => r.status === 400,
    });

    // Negative case, invalid payload
    currentTest = 'not valid id';
    res = testPostJson(route, { matchId: '123456789012345678901234' }, headers, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 404`]: (r) => r.status === 404,
    });
  }
  generateCatMatch(config, currentFeature, headers, {
    Authorization: `Bearer ${user2.accessToken}`,
  }, tags);

  currentTest = 'get all match cats';
  res = testGet(getRoute, {}, headers, tags);
  assert(res, currentFeature, config, {
    [`${currentTest} should return 200`]: (r) => r.status === 200,
  });
  let userCatMatch
  userCatMatch = res.json().data.find(
    /** @param {CatMatch} match */
    (match) => match.userCatDetail.hasMatched === false
      && match.issuedBy.email != user.email
      && match.matchCatDetail.hasMatched === false);

  currentTest = 'match cat approve'
  res = testPostJson(route, {
    matchId: userCatMatch.id,
  }, headers, tags);
  assert(res, currentFeature, config, {
    [`${currentTest} should return something`]: (r) => r.status
  }, {
    matchId: userCatMatch.id,
  });

  if (!config.POSITIVE_CASE) {
    currentTest = 'get all match cats';
    res = testGet(getRoute, {}, headers, tags);
    assert(res, currentFeature, config, {
      [`'${currentTest} should return 200`]: (r) => r.status === 200,
      [`'${currentTest} should not have approved match request`]: (r) => {
        try {
          return r.json().data.find((match) => match.id === userCatMatch.id) === undefined;
        } catch (e) {
          return false;
        }
      },
    }, {});
  }

  // if (!positivePayloadPassAssertTest) return null;
}


/**
 * Test manage cat match reject POST API
 * @param {Config} config
 * @param {User} user
 * @param {Object} tags
 */
export function TestPostManageCatReject(config, user, user2, tags = {}) {
  let res, currentTest;
  // eslint-disable-next-line no-undef
  const getRoute = `${__ENV.BASE_URL}/v1/cat/match`;
  // eslint-disable-next-line no-undef
  const route = `${__ENV.BASE_URL}/v1/cat/match/reject`;
  const currentFeature = `${TEST_NAME} | manage cat match reject`;
  if (!user) fail(`${currentFeature} fail due to user is empty`);

  const headers = {
    Authorization: `Bearer ${user.accessToken}`,
  };

  if (!config.POSITIVE_CASE) {
    // Negative case, no header
    currentTest = 'no header';
    res = testPostJson(route, {}, {}, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 401`]: (r) => r.status === 401,
    });

    // Negative case, no body
    currentTest = 'no body';
    res = testPostJson(route, {}, headers, tags, ['noContentType']);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 400`]: (r) => r.status === 400,
    });

    // Negative case, invalid payload
    currentTest = 'not valid id';
    res = testPostJson(route, { matchId: '123456789012345678901234' }, headers, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 404`]: (r) => r.status === 404,
    });
  }
  generateCatMatch(config, currentFeature, headers, {
    Authorization: `Bearer ${user2.accessToken}`,
  }, tags);

  currentTest = 'get all match cats';
  res = testGet(getRoute, {}, headers, tags);
  assert(res, currentFeature, config, {
    [`${currentTest} should return 200`]: (r) => r.status === 200,
  });

  let userCatMatch
  userCatMatch = res.json().data.find(
    /** @param {CatMatch} match */
    (match) => match.userCatDetail.hasMatched === false && match.matchCatDetail.hasMatched === false);

  currentTest = 'match cat reject'
  res = testPostJson(route, {
    matchId: userCatMatch.id,
  }, headers, tags);
  let positivePayloadPassAssertTest = assert(res, currentFeature, config, {
    [`${currentTest} should return something`]: (r) => r.status,
  });

  if (!config.POSITIVE_CASE) {
    currentTest = 'get all match cats';
    testGet(getRoute, {}, headers, tags);
    positivePayloadPassAssertTest = assert(res, currentFeature, config, {
      [`'${currentTest} should return 200`]: (r) => r.status === 200,
      [`'${currentTest} should not have rejected match request`]: (r) => {
        try {
          return r.json().data.find((match) => match.id === userCatMatch.id) === undefined;
        } catch (e) {
          return false;
        }
      },
    }, {});
  }

  if (!positivePayloadPassAssertTest) return null;
}
