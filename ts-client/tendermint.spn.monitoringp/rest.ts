/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface MonitoringpConnectionChannelID {
  channelID?: string;
}

export interface MonitoringpConsumerClientID {
  clientID?: string;
}

export interface MonitoringpMonitoringInfo {
  transmitted?: boolean;
  signatureCounts?: TypesSignatureCounts;
}

/**
 * Params defines the parameters for the module.
 */
export interface MonitoringpParams {
  /** @format int64 */
  lastBlockHeight?: string;
  consumerChainID?: string;
  consumerConsensusState?: TypesConsensusState;

  /** @format int64 */
  consumerUnbondingPeriod?: string;

  /** @format uint64 */
  consumerRevisionHeight?: string;
}

export interface MonitoringpQueryGetConnectionChannelIDResponse {
  ConnectionChannelID?: MonitoringpConnectionChannelID;
}

export interface MonitoringpQueryGetConsumerClientIDResponse {
  ConsumerClientID?: MonitoringpConsumerClientID;
}

export interface MonitoringpQueryGetMonitoringInfoResponse {
  MonitoringInfo?: MonitoringpMonitoringInfo;
}

/**
 * QueryParamsResponse is response type for the Query/Params RPC method.
 */
export interface MonitoringpQueryParamsResponse {
  /** Params defines the parameters for the module. */
  params?: MonitoringpParams;
}

export interface ProtobufAny {
  "@type"?: string;
}

export interface RpcStatus {
  /** @format int32 */
  code?: number;
  message?: string;
  details?: ProtobufAny[];
}

export interface TypesConsensusState {
  nextValidatorsHash?: string;
  timestamp?: string;
  root?: TypesMerkleRoot;
}

export interface TypesMerkleRoot {
  hash?: string;
}

export interface TypesSignatureCount {
  opAddress?: string;
  RelativeSignatures?: string;
}

export interface TypesSignatureCounts {
  /** @format uint64 */
  blockCount?: string;
  counts?: TypesSignatureCount[];
}

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, ResponseType } from "axios";

export type QueryParamsType = Record<string | number, any>;

export interface FullRequestParams extends Omit<AxiosRequestConfig, "data" | "params" | "url" | "responseType"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseType;
  /** request body */
  body?: unknown;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> extends Omit<AxiosRequestConfig, "data" | "cancelToken"> {
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
  secure?: boolean;
  format?: ResponseType;
}

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
}

export class HttpClient<SecurityDataType = unknown> {
  public instance: AxiosInstance;
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private secure?: boolean;
  private format?: ResponseType;

  constructor({ securityWorker, secure, format, ...axiosConfig }: ApiConfig<SecurityDataType> = {}) {
    this.instance = axios.create({ ...axiosConfig, baseURL: axiosConfig.baseURL || "" });
    this.secure = secure;
    this.format = format;
    this.securityWorker = securityWorker;
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  private mergeRequestParams(params1: AxiosRequestConfig, params2?: AxiosRequestConfig): AxiosRequestConfig {
    return {
      ...this.instance.defaults,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.instance.defaults.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  private createFormData(input: Record<string, unknown>): FormData {
    return Object.keys(input || {}).reduce((formData, key) => {
      const property = input[key];
      formData.append(
        key,
        property instanceof Blob
          ? property
          : typeof property === "object" && property !== null
          ? JSON.stringify(property)
          : `${property}`,
      );
      return formData;
    }, new FormData());
  }

  public request = async <T = any, _E = any>({
    secure,
    path,
    type,
    query,
    format,
    body,
    ...params
  }: FullRequestParams): Promise<AxiosResponse<T>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const responseFormat = (format && this.format) || void 0;

    if (type === ContentType.FormData && body && body !== null && typeof body === "object") {
      requestParams.headers.common = { Accept: "*/*" };
      requestParams.headers.post = {};
      requestParams.headers.put = {};

      body = this.createFormData(body as Record<string, unknown>);
    }

    return this.instance.request({
      ...requestParams,
      headers: {
        ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
        ...(requestParams.headers || {}),
      },
      params: query,
      responseType: responseFormat,
      data: body,
      url: path,
    });
  };
}

/**
 * @title monitoringp/connection_channel_id.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryConnectionChannelId
   * @summary Queries a ConnectionChannelID by index.
   * @request GET:/tendermint/spn/monitoringp/connection_channel_id
   */
  queryConnectionChannelId = (params: RequestParams = {}) =>
    this.request<MonitoringpQueryGetConnectionChannelIDResponse, RpcStatus>({
      path: `/tendermint/spn/monitoringp/connection_channel_id`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryConsumerClientId
   * @summary Queries a ConsumerClientID by index.
   * @request GET:/tendermint/spn/monitoringp/consumer_client_id
   */
  queryConsumerClientId = (params: RequestParams = {}) =>
    this.request<MonitoringpQueryGetConsumerClientIDResponse, RpcStatus>({
      path: `/tendermint/spn/monitoringp/consumer_client_id`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryMonitoringInfo
   * @summary Queries a MonitoringInfo by index.
   * @request GET:/tendermint/spn/monitoringp/monitoring_info
   */
  queryMonitoringInfo = (params: RequestParams = {}) =>
    this.request<MonitoringpQueryGetMonitoringInfoResponse, RpcStatus>({
      path: `/tendermint/spn/monitoringp/monitoring_info`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryParams
   * @summary Params queries the parameters of the module.
   * @request GET:/tendermint/spn/monitoringp/params
   */
  queryParams = (params: RequestParams = {}) =>
    this.request<MonitoringpQueryParamsResponse, RpcStatus>({
      path: `/tendermint/spn/monitoringp/params`,
      method: "GET",
      format: "json",
      ...params,
    });
}
