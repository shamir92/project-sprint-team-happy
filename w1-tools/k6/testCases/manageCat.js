import { fail } from 'k6';
import {
  generateRandomDescription, generateRandomImageUrl, generateRandomNumber, isValidDate, testDelete, testGet, testPutJson,

  generateTestObjects, generateUniqueName, isEqual, isExists, testPostJson, assert,
} from '../helper.js';
import { generateRandomCatBreed, generateRandomCatGender } from '../types/cat.js';

const manageCatNegativePayloads = generateTestObjects({
  name: { type: 'string', notNull: true, minLength: 1, maxLength: 30, },
  race: { type: 'string', notNull: true, enum: ['Persian', 'Maine Coon', 'Siamese', 'Ragdoll', 'Bengal', 'Sphynx', 'British Shorthair', 'Abyssinian', 'Scottish Fold', 'Birman'] },
  sex: { type: 'string', notNull: true, enum: ['male', 'female'] },
  ageInMonth: { type: 'number', notNull: true, min: 1, max: 120082, },
  description: { type: 'string', notNull: true, minLength: 1, maxLength: 200, },
  imageUrls: { type: 'array', notNull: true, items: { type: 'string', isUrl: true } },
}, {
  name: generateUniqueName(),
  race: 'Persian',
  ageInMonth: 12,
  sex: 'female',
  description: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
  imageUrls: [generateRandomImageUrl()],
});

const TEST_NAME = '(manage cat)';
/**
 * Test manage cat POST API
 * @param {Config} config
 * @param {Object} tags
 * @param {User} user
 * @param {number} ageInMonth
 * @return {Cat} cat
 */
export function TestPostManageCat(config, user, tags = {}, ageInMonth = 0) {
  let res, currentTest;
  // eslint-disable-next-line no-undef
  const route = `${__ENV.BASE_URL}/v1/cat`;
  const currentFeature = `${TEST_NAME} | post manage cat`;
  if (!user) fail(`${currentFeature} fail due to user is empty`);

  const catBreed = generateRandomCatBreed();
  const catSex = generateRandomCatGender();
  const positivePayload = {
    name: generateUniqueName(),
    race: catBreed,
    ageInMonth: ageInMonth || generateRandomNumber(1, 120082),
    sex: catSex,
    description: generateRandomDescription(200),
    imageUrls: [generateRandomImageUrl()],
  };
  const headers = {
    Authorization: `Bearer ${user.accessToken}`,
  };

  if (!config.POSITIVE_CASE) {
    currentTest = 'no header';
    res = testPostJson(route, {}, {}, tags, ['noContentType']);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 401`]: (r) => r.status === 401,
    });


    currentTest = 'no payload';
    res = testPostJson(route, {}, headers, tags, ['noContentType']);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 400`]: (r) => r.status === 400,
    });


    manageCatNegativePayloads.forEach((payload) => {
      currentTest = 'invalid payload';
      res = testPostJson(route, payload, headers, tags);
      assert(res, currentFeature, config, {
        [`${currentTest} should return 400`]: (r) => r.status === 400,
      }, payload);
    });
  }


  currentTest = 'valid payload';
  res = testPostJson(route, positivePayload, headers, tags);
  const positivePayloadPassAssertTest = assert(res, currentFeature, config, {
    [`${currentTest} should return 201`]: (r) => r.status === 201,
    [`${currentTest} should have id`]: (r) => isExists(r, 'data.id'),
    [`${currentTest} should have createdAt`]: (r) => isExists(r, 'data.createdAt'),
    [`${currentTest} createdAt should be in ISO 8601 format`]: (r) => isValidDate(r.json().data.createdAt),
  }, positivePayload);


  if (!positivePayloadPassAssertTest) return null;

  return Object.assign(positivePayload, res.json().data);

}

/**
 * Test manage cat GET API
 * @param {Config} config
 * @param {Object} tags
 * @param {User} user
 * @param {Cat} cat
 */
export function TestGetManageCat(config, user, cat, tags = {}) {
  let res, currentTest;
  // eslint-disable-next-line no-undef
  const route = `${__ENV.BASE_URL}/v1/cat`;
  const currentFeature = `${TEST_NAME} | get manage cat`;
  if (!user) fail(`${currentFeature} fail due to user is empty`);
  if (!cat) fail(`${currentFeature} fail due to cat is empty`);

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

  // Positive case, get all cat
  const commonChecks = (currFeat) => ({
    [`${currFeat} should return 200`]: (r) => r.status === 200,
    [`${currFeat} should have all id`]: (r) => isExists(r, 'data.[].id'),
    [`${currFeat} should have all name`]: (r) => isExists(r, 'data.[].name'),
    [`${currFeat} should have all race`]: (r) => isExists(r, 'data.[].race'),
    [`${currFeat} should have all sex`]: (r) => isExists(r, 'data.[].sex'),
    [`${currFeat} should have all ageInMonth`]: (r) => isExists(r, 'data.[].ageInMonth'),
    [`${currFeat} should have all imageUrls`]: (r) => isExists(r, 'data.[].imageUrls'),
    [`${currFeat} should have all description`]: (r) => isExists(r, 'data.[].description'),
    [`${currFeat} should have all hasMatched`]: (r) => isExists(r, 'data.[].hasMatched'),
    [`${currFeat} should have all createdAt`]: (r) => isExists(r, 'data.[].createdAt'),
  });
  currentTest = 'get all cat';
  res = testGet(route, {}, headers, tags);
  const positivePayloadPassAssertTest = assert(res, currentFeature, config, commonChecks(currentTest), {});

  if (!positivePayloadPassAssertTest) return null;

  currentTest = 'search by id';
  res = testGet(route, { id: cat.id }, headers, tags);
  assert(res, currentFeature, config, Object.assign(commonChecks(currentTest), {
    [`${currentTest} should have the value equal to id query param`]: (r) => isEqual(r, 'data.[].id', cat.id),
  }), { id: cat.id });

  currentTest = 'search by name';
  res = testGet(route, { search: 'a' }, headers, tags);
  assert(res, currentFeature, config, Object.assign(commonChecks(currentTest), {
    [`${currentTest} should have the value equal to search query param`]: (r) => {
      try {
        return r.json().data.every((c) => c.name.toLowerCase().includes('a'));
      } catch (e) {
        return false;
      }
    },
  }), { search: 'a' });


  currentTest = 'search by cat race'
  res = testGet(route, { race: cat.race }, headers, tags);
  assert(res, currentFeature, config, Object.assign(commonChecks(currentTest), {
    'valid payload should have the value equal to race query param': (r) => isEqual(r, 'data.[].race', cat.race),
  }), { race: cat.race });


  currentTest = 'search by gender'
  res = testGet(route, { sex: cat.sex }, headers, tags);
  assert(res, currentFeature, config, Object.assign(commonChecks(currentTest), {
    'valid payload should have the value equal to gender query param': (r) => isEqual(r, 'data.[].sex', cat.sex),
  }), { sex: cat.gender });


  currentTest = 'search by owned'
  res = testGet(route, { owned: true }, headers, tags);
  assert(res, currentFeature, config, Object.assign(commonChecks(currentTest), {
    'valid payload should have the value equal to owned query param': (r) => isEqual(r, 'data.[].id', cat.id),
  }), { owned: true });


  const middleAgeInMonth = generateRandomNumber(3, 120080);
  if (!config.GACHA) {
    const positivePayload = (ageInMonth, gender) => ({
      name: generateUniqueName(),
      race: generateRandomCatBreed(),
      ageInMonth: ageInMonth,
      sex: gender,
      description: generateRandomDescription(200),
      imageUrls: [generateRandomImageUrl()],
    });
    currentTest = 'create male cat with ageInMonth less than middleAgeInMonth';
    const lessThanAgePayloadMale = positivePayload(generateRandomNumber(1, middleAgeInMonth - 1), "male");
    res = testPostJson(route, lessThanAgePayloadMale, headers, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 201`]: (r) => r.status === 201,
    }, lessThanAgePayloadMale);


    currentTest = 'create male cat with ageInMonth more than middleAgeInMonth';
    const moreThanAgePayloadMale = positivePayload(generateRandomNumber(middleAgeInMonth + 1, 120082), "male");
    res = testPostJson(route, moreThanAgePayloadMale, headers, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 201`]: (r) => r.status === 201,
    }, moreThanAgePayloadMale);


    currentTest = 'create male cat with ageInMonth equal to middleAgeInMonth';
    const centerAgePayloadMale = positivePayload(middleAgeInMonth, "male");
    res = testPostJson(route, centerAgePayloadMale, headers, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 201`]: (r) => r.status === 201,
    }, centerAgePayloadMale);

    currentTest = 'create female cat with ageInMonth less than middleAgeInMonth';
    const lessThanAgePayloadFemale = positivePayload(generateRandomNumber(1, middleAgeInMonth - 1), "female");
    res = testPostJson(route, lessThanAgePayloadFemale, headers, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 201`]: (r) => r.status === 201,
    }, lessThanAgePayloadFemale);


    currentTest = 'create female cat with ageInMonth more than middleAgeInMonth';
    const moreThanAgePayloadFemale = positivePayload(generateRandomNumber(middleAgeInMonth + 1, 120082), "female");
    res = testPostJson(route, moreThanAgePayloadFemale, headers, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 201`]: (r) => r.status === 201,
    }, moreThanAgePayloadFemale);


    currentTest = 'create female cat with ageInMonth equal to middleAgeInMonth';
    const centerAgePayloadFemale = positivePayload(middleAgeInMonth, "female");
    res = testPostJson(route, centerAgePayloadFemale, headers, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 201`]: (r) => r.status === 201,
    }, centerAgePayloadFemale);
  }

  currentTest = 'search by less than ageInMonth';
  res = testGet(route, { ageInMonth: `<${middleAgeInMonth}` }, headers, tags);
  assert(res, currentFeature, config, Object.assign(commonChecks(currentTest), {
    [`${currentTest} should have the value less than in ageInMonth param`]: (r) => {
      try {
        return r.json().data.every((c) => c.ageInMonth < middleAgeInMonth);
      } catch (e) {
        return false;
      }
    },
  }), { ageInMonth: `<${middleAgeInMonth}` });


  currentTest = 'search by more than ageInMonth';
  res = testGet(route, { ageInMonth: `>${middleAgeInMonth}` }, headers, tags);
  assert(res, currentFeature, config, Object.assign(commonChecks(currentTest), {
    [`${currentTest} should have the value more than in ageInMonth param`]: (r) => {
      try {
        return r.json().data.every((c) => c.ageInMonth > middleAgeInMonth);
      } catch (e) {
        return false;
      }
    },
  }), { ageInMonth: `>${middleAgeInMonth}` });

  currentTest = 'search by equal to ageInMonth';
  res = testGet(route, { ageInMonth: `=${middleAgeInMonth}` }, headers, tags);
  assert(res, currentFeature, config, Object.assign(commonChecks(currentTest), {
    [`${currentTest} should have the value equal to ageInMonth param`]: (r) => {
      try {
        return r.json().data.every((c) => c.ageInMonth === middleAgeInMonth);
      } catch (e) {
        return false;
      }
    },
  }), { ageInMonth: `=${middleAgeInMonth}` });
}

/**
 * Test manage cat PUT API
 * @param {Config} config
 * @param {User} user
 * @param {Object} tags
 * @return {import("../types/cat.js").Cat} cat
 */
export function TestPutManageCat(config, user, tags = {}) {
  let res, currentTest;
  // eslint-disable-next-line no-undef
  const getRoute = `${__ENV.BASE_URL}/v1/cat`;
  const currentFeature = `${TEST_NAME} | put manage cat`;
  if (!user) fail(`${currentFeature} fail due to user is empty`);

  const headers = {
    Authorization: `Bearer ${user.accessToken}`,
  };

  currentTest = 'search by owned'
  res = testGet(getRoute, { owned: true }, headers, tags);
  assert(res, currentFeature, config, {
    [`${currentTest} should return 200`]: (r) => r.status === 200,
  }, { owned: true });
  let cat = res.json().data[generateRandomNumber(0, res.json().data.length - 1)];

  // eslint-disable-next-line no-undef
  const route = `${__ENV.BASE_URL}/v1/cat/${cat.id}`;
  const positivePayload = {
    name: generateUniqueName(),
    race: generateRandomCatBreed(),
    ageInMonth: generateRandomNumber(1, 120082),
    sex: generateRandomCatGender(cat.sex),
    description: generateRandomDescription(200),
    imageUrls: [generateRandomImageUrl()],
  };

  if (!config.POSITIVE_CASE) {
    currentTest = 'no header';
    res = testPutJson(route, {}, {}, tags, ['noContentType']);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 401`]: (r) => r.status === 401,
    });


    currentTest = 'no payload';
    res = testPutJson(route, {}, headers, tags, ['noContentType']);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 401`]: (r) => r.status === 400,
    });


    manageCatNegativePayloads.forEach((payload) => {
      currentTest = 'invalid payload';
      res = testPutJson(route, payload, headers, tags);
      assert(res, currentFeature, config, {
        [`${currentTest} should return 400`]: (r) => r.status === 400,
      }, payload);
    });
  }

  currentTest = 'update cat';
  res = testPutJson(route, positivePayload, headers, tags);
  let positivePayloadPassAssertTest = assert(res, currentFeature, config, {
    [`${currentTest} should return 200`]: (r) => r.status === 200,
  }, positivePayload);


  currentTest = 'updated cat should be the same as the one in the get route';
  res = testGet(getRoute, { id: cat.id }, headers, tags);
  positivePayloadPassAssertTest = assert(res, currentFeature, config, {
    [`${currentTest} should return 200`]: (r) => r.status === 200,
    [`${currentTest} should have the value equal to id query param`]: (r) => isEqual(r, 'data.[].id', cat.id),
    [`${currentTest} should have the value equal to name in payload`]: (r) => isEqual(r, 'data.[].name', positivePayload.name),
    [`${currentTest} should have the value equal to race in payload`]: (r) => isEqual(r, 'data.[].race', positivePayload.race),
    [`${currentTest} should have the value equal to ageInMonth in payload`]: (r) => isEqual(r, 'data.[].ageInMonth', positivePayload.ageInMonth),
    [`${currentTest} should have the value equal to sex in payload`]: (r) => isEqual(r, 'data.[].sex', positivePayload.sex),
    [`${currentTest} should have the value equal to description in payload`]: (r) => isEqual(r, 'data.[].description', positivePayload.description),
  }, { id: cat.id });


  if (!positivePayloadPassAssertTest) return null;

  return res.json().data;
}

/**
 * Test manage cat DELETE API
 * @param {Config} config
 * @param {User} user
 * @param {Object} tags
 * @return {import("../types/cat.js").Cat} cat
 */
export function TestDeleteManageCat(config, user, tags = {}) {
  let res, currentTest;
  // eslint-disable-next-line no-undef
  const getRoute = `${__ENV.BASE_URL}/v1/cat`;
  const currentFeature = `${TEST_NAME} | delete manage cat`;
  if (!user) fail(`${currentFeature} fail due to user is empty`);

  const headers = {
    Authorization: `Bearer ${user.accessToken}`,
  };

  currentTest = 'search by owned'
  res = testGet(getRoute, { owned: true }, headers, tags);
  assert(res, currentFeature, config, {
    [`${currentTest} should return 200`]: (r) => r.status === 200,
  }, { owned: true });
  let cat = res.json().data[generateRandomNumber(0, res.json().data.length - 1)];

  // eslint-disable-next-line no-undef
  const route = `${__ENV.BASE_URL}/v1/cat/${cat.id}`;

  if (!config.POSITIVE_CASE) {
    currentTest = 'no header';
    res = testDelete(route, {}, {}, tags);
    assert(res, currentFeature, config, {
      [`${currentTest} should return 401`]: (r) => r.status === 401,
    });
  }

  currentTest = 'delete cat';
  res = testDelete(route, {}, headers, tags);
  let positivePayloadPassAssertTest = assert(res, currentFeature, config, {
    [`${currentTest} should return 200`]: (r) => r.status === 200,
  });

  currentTest = 'deleted cat should not be found in the get route';
  res = testGet(getRoute, { id: cat.id }, headers, tags);
  positivePayloadPassAssertTest = assert(res, currentFeature, config, {
    [`${currentTest} should return 200`]: (r) => r.status === 200,
    [`${currentTest} should return empty array`]: (r) => {
      try {
        return r.json().data.length === 0;
      } catch (e) {
        return false;
      }
    },
  }, { id: cat.id });

  if (!positivePayloadPassAssertTest) return null;

  return res.json().data;
}
