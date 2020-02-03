# gopu

该服务是一个基于gin框架的http账户管理服务的简单实现，功能包括
* JWT令牌验证
* Role权限管理（使用Casbin）
* 用户注册/登录
* 头像上传
* 邮箱验证
* ...

对于注册的用户，邮箱的验证使用验证码进行验证，即需要通过邮箱验证码才能完成注册，防止恶意邮箱注册。

API包含：
1. 登录验证
   * POST   /v1/session :用户登录
   * DELETE /v1/session :用户登出
   * GET /v1/session/refresh_token :令牌刷新
2. 用户信息
   * POST /v1/user/password/reset_code :发送用户密码重置码到邮箱
   * POST /v1/user/register_code :发送用户注册邮箱验证码到邮箱
   * POST /v1/user :注册用户
   * GET  /v1/user/:id :获取对应用户id的用户信息
   * PUT  /v1/user/:id/password :设置用户id对应的密码信息
   * GET  /v1/current_user :根据登录令牌获取当前用户信息
   * PUT  /v1/user/:id/profile :设置对应用户id的数据信息
   * DELETE /v1/user/:id :删除对应用户id的用户信息
   * GET  /v1/user :获取用户信息列表
3. 角色管理
   * POST /v1/role :创建一个角色
   * DELETE /v1/role/:name :删除对应角色名称name的角色信息
   * GET /v1/role/:name/user :获取对应角色名称name的所有用户
   * GET /v1/role/:name :获取对应角色名称name的角色信息
   * GET /v1/role :获取角色列表
   * POST /v1/role/:name/user/:id 添加用户id到角色name列表中
   * DELETE /v1/role/:name/user/:id 从角色name列表中删除用户id
