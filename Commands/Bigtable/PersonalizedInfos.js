'use strict';
/*                          _       _
 *__      _____  __ ___   ___  __ _| |_ ___
 *\ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
 * \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
 *  \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
 *
 * Copyright © 2016 Weaviate. All rights reserved.
 * LICENSE: https://github.com/weaviate/weaviate/blob/master/LICENSE
 * AUTHOR: Bob van Luijt (bob@weaviate.com)
 * See www.weaviate.com for details
 * See package.json for author and maintainer info
 * Contact: @weaviate_iot / yourfriends@weaviate.com
 */

/** Class Commands_PersonalizedInfos */
module.exports = class Commands_PersonalizedInfos { // Class: Commands_{resources.className}

    /**
     * Constructor for this Command
     * @param {object} req  - The request
     * @param {object} res  - The response
     * @param {object} next - Next() function
     */
    constructor(req, res, next) {
        this.req  = req;
        this.res  = res;
        this.next = next;
    }

    /**
     * Returns the personalized info for device.
     * @param {object} commandAttributes  - All attributes needed to exec the command
     * @return {promise} Returns a promise with the correct object
     */
    $Get(commandAttributes) {
        return new Promise((resolve, reject) => {
            // resolve with kind and token
            resolve({
                id: 'me',
                kind: 'weave#personalizedInfo',
                lastUseTimeMs: 12345,
                location: 'Study',
                name: 'This is my name'
            });
        });
    }

    /**
     * Update the personalized info for device. This method supports patch semantics.
     * @param {object} commandAttributes  - All attributes needed to exec the command
     * @return {promise} Returns a promise with the correct object
     */
    $Patch(commandAttributes) {
        return new Promise((resolve, reject) => {
            // resolve with kind and token
            resolve({
                id: 'me',
                kind: 'weave#personalizedInfo',
                lastUseTimeMs: 12345,
                location: 'Study',
                name: 'This is my name'
            });
        });
    }

    /**
     * Update the personalized info for device.
     * @param {object} commandAttributes  - All attributes needed to exec the command
     * @return {promise} Returns a promise with the correct object
     */
    $Update(commandAttributes) {
        return new Promise((resolve, reject) => {
            // resolve with kind and token
            resolve({
                id: 'me',
                kind: 'weave#personalizedInfo',
                lastUseTimeMs: 12345,
                location: 'Study',
                name: 'This is my name'
            });
        });
    }

};
