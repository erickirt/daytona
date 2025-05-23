/* tslint:disable */

/**
 * Daytona
 * Daytona AI platform API Docs
 *
 * The version of the OpenAPI document: 1.0
 * Contact: support@daytona.com
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

/**
 *
 * @export
 * @interface VolumeDto
 */
export interface VolumeDto {
  /**
   * Volume ID
   * @type {string}
   * @memberof VolumeDto
   */
  id: string
  /**
   * Volume name
   * @type {string}
   * @memberof VolumeDto
   */
  name: string
  /**
   * Organization ID
   * @type {string}
   * @memberof VolumeDto
   */
  organizationId: string
  /**
   * Volume state
   * @type {string}
   * @memberof VolumeDto
   */
  state: VolumeDtoStateEnum
  /**
   * Creation timestamp
   * @type {string}
   * @memberof VolumeDto
   */
  createdAt: string
  /**
   * Last update timestamp
   * @type {string}
   * @memberof VolumeDto
   */
  updatedAt: string
  /**
   * Last used timestamp
   * @type {string}
   * @memberof VolumeDto
   */
  lastUsedAt: string
}

export const VolumeDtoStateEnum = {
  CREATING: 'creating',
  READY: 'ready',
  PENDING_CREATE: 'pending_create',
  PENDING_DELETE: 'pending_delete',
  DELETING: 'deleting',
  DELETED: 'deleted',
  ERROR: 'error',
} as const

export type VolumeDtoStateEnum = (typeof VolumeDtoStateEnum)[keyof typeof VolumeDtoStateEnum]
