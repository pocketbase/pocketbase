declare abstract class BaseModel {
    id: string;
    created: string;
    updated: string;
    constructor(data?: {
        [key: string]: any;
    });
    /**
     * Loads `data` into the current model.
     */
    load(data: {
        [key: string]: any;
    }): void;
    /**
     * Returns whether the current loaded data represent a stored db record.
     */
    get isNew(): boolean;
    /**
     * Robust deep clone of a model.
     */
    clone(): BaseModel;
    /**
     * Exports all model properties as a new plain object.
     */
    export(): {
        [key: string]: any;
    };
}
declare class Record extends BaseModel {
    [key: string]: any;
    "@collectionId": string;
    "@collectionName": string;
    "@expand": {
        [key: string]: any;
    };
    /**
     * @inheritdoc
     */
    load(data: {
        [key: string]: any;
    }): void;
}
declare class User extends BaseModel {
    email: string;
    verified: boolean;
    lastResetSentAt: string;
    lastVerificationSentAt: string;
    profile: null | Record;
    /**
     * @inheritdoc
     */
    load(data: {
        [key: string]: any;
    }): void;
}
declare class Admin extends BaseModel {
    avatar: number;
    email: string;
    lastResetSentAt: string;
    /**
     * @inheritdoc
     */
    load(data: {
        [key: string]: any;
    }): void;
}
type AuthStore = {
    /**
     * Retrieves the stored token (if any).
     */
    readonly token: string;
    /**
     * Retrieves the stored model data (if any).
     */
    readonly model: User | Admin | {};
    /**
     * Checks if the store has valid (aka. existing and unexpired) token.
     */
    readonly isValid: boolean;
    /**
     * Saves new token and model data in the auth store.
     */
    save(token: string, model: User | Admin | {}): void;
    /**
     * Removes the stored token and model data form the auth store.
     */
    clear(): void;
};
/**
 * BaseService class that should be inherited from all API services.
 */
declare abstract class BaseService {
    readonly client: Client;
    constructor(client: Client);
}
declare class Settings extends BaseService {
    /**
     * Fetch all available app settings.
     */
    getAll(queryParams?: {}): Promise<{
        [key: string]: any;
    }>;
    /**
     * Bulk updates app settings.
     */
    update(bodyParams?: {}, queryParams?: {}): Promise<{
        [key: string]: any;
    }>;
}
declare class ListResult<M extends BaseModel> {
    page: number;
    perPage: number;
    totalItems: number;
    items: Array<M>;
    constructor(page: number, perPage: number, totalItems: number, items: Array<M>);
}
declare abstract class BaseCrudService<M extends BaseModel> extends BaseService {
    /**
     * Response data decoder.
     */
    abstract decode(data: {
        [key: string]: any;
    }): M;
    /**
     * Returns a promise with all list items batch fetched at once.
     */
    protected _getFullList(basePath: string, batchSize?: number, queryParams?: {}): Promise<Array<M>>;
    /**
     * Returns paginated items list.
     */
    protected _getList(basePath: string, page?: number, perPage?: number, queryParams?: {}): Promise<ListResult<M>>;
    /**
     * Returns single item by its id.
     */
    protected _getOne(basePath: string, id: string, queryParams?: {}): Promise<M>;
    /**
     * Creates a new item.
     */
    protected _create(basePath: string, bodyParams?: {}, queryParams?: {}): Promise<M>;
    /**
     * Updates an existing item by its id.
     */
    protected _update(basePath: string, id: string, bodyParams?: {}, queryParams?: {}): Promise<M>;
    /**
     * Deletes an existing item by its id.
     */
    protected _delete(basePath: string, id: string, bodyParams?: {}, queryParams?: {}): Promise<boolean>;
}
declare abstract class CrudService<M extends BaseModel> extends BaseCrudService<M> {
    /**
     * Base path for the crud actions (without trailing slash, eg. '/admins').
     */
    abstract baseCrudPath(): string;
    /**
     * Returns a promise with all list items batch fetched at once.
     */
    getFullList(batchSize?: number, queryParams?: {}): Promise<Array<M>>;
    /**
     * Returns paginated items list.
     */
    getList(page?: number, perPage?: number, queryParams?: {}): Promise<ListResult<M>>;
    /**
     * Returns single item by its id.
     */
    getOne(id: string, queryParams?: {}): Promise<M>;
    /**
     * Creates a new item.
     */
    create(bodyParams?: {}, queryParams?: {}): Promise<M>;
    /**
     * Updates an existing item by its id.
     */
    update(id: string, bodyParams?: {}, queryParams?: {}): Promise<M>;
    /**
     * Deletes an existing item by its id.
     */
    delete(id: string, bodyParams?: {}, queryParams?: {}): Promise<boolean>;
}
type AdminAuthResponse = {
    [key: string]: any;
    token: string;
    admin: Admin;
};
declare class Admins extends CrudService<Admin> {
    /**
     * @inheritdoc
     */
    decode(data: {
        [key: string]: any;
    }): Admin;
    /**
     * @inheritdoc
     */
    baseCrudPath(): string;
    /**
     * Prepare successfull authorize response.
     */
    protected authResponse(responseData: any): AdminAuthResponse;
    /**
     * Authenticate an admin account by its email and password
     * and returns a new admin token and data.
     *
     * On success this method automatically updates the client's AuthStore data.
     */
    authViaEmail(email: string, password: string, bodyParams?: {}, queryParams?: {}): Promise<AdminAuthResponse>;
    /**
     * Refreshes the current admin authenticated instance and
     * returns a new token and admin data.
     *
     * On success this method automatically updates the client's AuthStore data.
     */
    refresh(bodyParams?: {}, queryParams?: {}): Promise<AdminAuthResponse>;
    /**
     * Sends admin password reset request.
     */
    requestPasswordReset(email: string, bodyParams?: {}, queryParams?: {}): Promise<boolean>;
    /**
     * Confirms admin password reset request.
     */
    confirmPasswordReset(passwordResetToken: string, password: string, passwordConfirm: string, bodyParams?: {}, queryParams?: {}): Promise<AdminAuthResponse>;
}
type UserAuthResponse = {
    [key: string]: any;
    token: string;
    user: User;
};
type AuthProviderInfo = {
    name: string;
    state: string;
    codeVerifier: string;
    codeChallenge: string;
    codeChallengeMethod: string;
    authUrl: string;
};
type AuthMethodsList = {
    [key: string]: any;
    emailPassword: boolean;
    authProviders: Array<AuthProviderInfo>;
};
declare class Users extends CrudService<User> {
    /**
     * @inheritdoc
     */
    decode(data: {
        [key: string]: any;
    }): User;
    /**
     * @inheritdoc
     */
    baseCrudPath(): string;
    /**
     * Prepare successfull authorization response.
     */
    protected authResponse(responseData: any): UserAuthResponse;
    /**
     * Returns all available application auth methods.
     */
    listAuthMethods(queryParams?: {}): Promise<AuthMethodsList>;
    /**
     * Authenticate a user via its email and password.
     *
     * On success, this method also automatically updates
     * the client's AuthStore data and returns:
     * - new user authentication token
     * - the authenticated user model record
     */
    authViaEmail(email: string, password: string, bodyParams?: {}, queryParams?: {}): Promise<UserAuthResponse>;
    /**
     * Authenticate a user via OAuth2 client provider.
     *
     * On success, this method also automatically updates
     * the client's AuthStore data and returns:
     * - new user authentication token
     * - the authenticated user model record
     * - the OAuth2 user profile data (eg. name, email, avatar, etc.)
     */
    authViaOAuth2(provider: string, code: string, codeVerifier: string, redirectUrl: string, bodyParams?: {}, queryParams?: {}): Promise<UserAuthResponse>;
    /**
     * Refreshes the current user authenticated instance and
     * returns a new token and user data.
     *
     * On success this method also automatically updates the client's AuthStore data.
     */
    refresh(bodyParams?: {}, queryParams?: {}): Promise<UserAuthResponse>;
    /**
     * Sends user password reset request.
     */
    requestPasswordReset(email: string, bodyParams?: {}, queryParams?: {}): Promise<boolean>;
    /**
     * Confirms user password reset request.
     */
    confirmPasswordReset(passwordResetToken: string, password: string, passwordConfirm: string, bodyParams?: {}, queryParams?: {}): Promise<UserAuthResponse>;
    /**
     * Sends user verification email request.
     */
    requestVerification(email: string, bodyParams?: {}, queryParams?: {}): Promise<boolean>;
    /**
     * Confirms user email verification request.
     */
    confirmVerification(verificationToken: string, bodyParams?: {}, queryParams?: {}): Promise<UserAuthResponse>;
    /**
     * Sends an email change request to the authenticated user.
     */
    requestEmailChange(newEmail: string, bodyParams?: {}, queryParams?: {}): Promise<boolean>;
    /**
     * Confirms user new email address.
     */
    confirmEmailChange(emailChangeToken: string, password: string, bodyParams?: {}, queryParams?: {}): Promise<UserAuthResponse>;
}
declare class SchemaField {
    id: string;
    name: string;
    type: string;
    system: boolean;
    required: boolean;
    unique: boolean;
    options: {
        [key: string]: any;
    };
    constructor(data?: {
        [key: string]: any;
    });
    /**
     * Loads `data` into the field.
     */
    load(data: {
        [key: string]: any;
    }): void;
}
declare class Collection extends BaseModel {
    name: string;
    schema: Array<SchemaField>;
    system: boolean;
    listRule: null | string;
    viewRule: null | string;
    createRule: null | string;
    updateRule: null | string;
    deleteRule: null | string;
    /**
     * @inheritdoc
     */
    load(data: {
        [key: string]: any;
    }): void;
}
declare class Collections extends CrudService<Collection> {
    /**
     * @inheritdoc
     */
    decode(data: {
        [key: string]: any;
    }): Collection;
    /**
     * @inheritdoc
     */
    baseCrudPath(): string;
}
declare abstract class SubCrudService<M extends BaseModel> extends BaseCrudService<M> {
    /**
     * Base path for the crud actions (without trailing slash, eg. '/collections/{:sub}/records').
     */
    abstract baseCrudPath(sub: string): string;
    /**
     * Returns a promise with all list items batch fetched at once.
     */
    getFullList(sub: string, batchSize?: number, queryParams?: {}): Promise<Array<M>>;
    /**
     * Returns paginated items list.
     */
    getList(sub: string, page?: number, perPage?: number, queryParams?: {}): Promise<ListResult<M>>;
    /**
     * Returns single item by its id.
     */
    getOne(sub: string, id: string, queryParams?: {}): Promise<M>;
    /**
     * Creates a new item.
     */
    create(sub: string, bodyParams?: {}, queryParams?: {}): Promise<M>;
    /**
     * Updates an existing item by its id.
     */
    update(sub: string, id: string, bodyParams?: {}, queryParams?: {}): Promise<M>;
    /**
     * Deletes an existing item by its id.
     */
    delete(sub: string, id: string, bodyParams?: {}, queryParams?: {}): Promise<boolean>;
}
declare class Records extends SubCrudService<Record> {
    /**
     * @inheritdoc
     */
    decode(data: {
        [key: string]: any;
    }): Record;
    /**
     * @inheritdoc
     */
    baseCrudPath(collectionIdOrName: string): string;
    /**
     * Builds and returns an absolute record file url.
     */
    getFileUrl(record: Record, filename: string, queryParams?: {}): string;
}
declare class LogRequest extends BaseModel {
    url: string;
    method: string;
    status: number;
    auth: string;
    ip: string;
    referer: string;
    userAgent: string;
    meta: null | {
        [key: string]: any;
    };
    /**
     * @inheritdoc
     */
    load(data: {
        [key: string]: any;
    }): void;
}
type HourlyStats = {
    total: number;
    date: string;
};
declare class Logs extends BaseService {
    /**
     * Returns paginated logged requests list.
     */
    getRequestsList(page?: number, perPage?: number, queryParams?: {}): Promise<ListResult<LogRequest>>;
    /**
     * Returns a single logged request by its id.
     */
    getRequest(id: string, queryParams?: {}): Promise<LogRequest>;
    /**
     * Returns request logs statistics.
     */
    getRequestsStats(queryParams?: {}): Promise<Array<HourlyStats>>;
}
interface MessageData {
    [key: string]: any;
    action: string;
    record: Record;
}
interface SubscriptionFunc {
    (data: MessageData): void;
}
declare class Realtime extends BaseService {
    private clientId;
    private eventSource;
    private subscriptions;
    /**
     * Inits the sse connection (if not already) and register the subscription.
     */
    subscribe(subscription: string, callback: SubscriptionFunc): Promise<void>;
    /**
     * Unsubscribe from a subscription.
     *
     * If the `subscription` argument is not set,
     * then the client will unsubscibe from all registered subscriptions.
     *
     * The related sse connection will be autoclosed if after the
     * unsubscribe operations there are no active subscriptions left.
     */
    unsubscribe(subscription?: string): Promise<void>;
    private submitSubscriptions;
    private addSubscriptionListeners;
    private removeSubscriptionListeners;
    private connectHandler;
    private connect;
    private disconnect;
}
/**
 * PocketBase JS Client.
 */
declare class Client {
    baseUrl: string;
    lang: string;
    AuthStore: AuthStore;
    readonly Settings: Settings;
    readonly Admins: Admins;
    readonly Users: Users;
    readonly Collections: Collections;
    readonly Records: Records;
    readonly Logs: Logs;
    readonly Realtime: Realtime;
    private cancelControllers;
    constructor(baseUrl?: string, lang?: string, authStore?: AuthStore | null);
    /**
     * Cancels single request by its cancellation key.
     */
    cancelRequest(cancelKey: string): Client;
    /**
     * Cancels all pending requests.
     */
    cancelAllRequests(): Client;
    /**
     * Sends an api http request.
     */
    send(path: string, reqConfig: {
        [key: string]: any;
    }): Promise<any>;
    /**
     * Builds a full client url by safely concatenating the provided path.
     */
    buildUrl(path: string): string;
    /**
     * Serializes the provided query parameters into a query string.
     */
    private serializeQueryParams;
}
/**
 * ClientResponseError is a custom Error class that is intended to wrap
 * and normalize any error thrown by `Client.send()`.
 */
declare class ClientResponseError extends Error {
    url: string;
    status: number;
    data: {
        [key: string]: any;
    };
    isAbort: boolean;
    originalError: any;
    constructor(errData?: any);
    // Make a POJO's copy of the current error class instance.
    // @see https://github.com/vuex-orm/vuex-orm/issues/255
    toJSON(): this;
}
/**
 * Default token store for browsers with auto fallback
 * to runtime/memory if local storage is undefined (eg. node env).
 */
declare class LocalAuthStore implements AuthStore {
    private fallback;
    private storageKey;
    constructor(storageKey?: string);
    /**
     * @inheritdoc
     */
    get token(): string;
    /**
     * @inheritdoc
     */
    get model(): User | Admin | {};
    /**
     * @inheritdoc
     */
    get isValid(): boolean;
    /**
     * @inheritdoc
     */
    save(token: string, model: User | Admin | {}): void;
    /**
     * @inheritdoc
     */
    clear(): void;
    // ---------------------------------------------------------------
    // Internal helpers:
    // ---------------------------------------------------------------
    /**
     * Retrieves `key` from the browser's local storage
     * (or runtime/memory if local storage is undefined).
     */
    private _storageGet;
    /**
     * Stores a new data in the browser's local storage
     * (or runtime/memory if local storage is undefined).
     */
    private _storageSet;
    /**
     * Removes `key` from the browser's local storage and the runtime/memory.
     */
    private _storageRemove;
}
export { Client as default, ClientResponseError, LocalAuthStore, User, Admin, Collection, Record, LogRequest, SchemaField };
