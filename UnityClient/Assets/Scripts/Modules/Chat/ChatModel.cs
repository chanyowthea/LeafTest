using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using Msg;
using UnityEngine;
using Util;
using Net;

public class ChatModel : BaseModel<ChatModel>
{
    protected override void InitAddTocHandler()
    {
        AddTocHandler(typeof(TocChat), STocChat);
        AddTocHandler(typeof(LoginRes_0X0101), LoginRes);
        AddTocHandler(typeof(LoginFaild), LoginFailed);
        AddTocHandler(typeof(LoginSuccessfull), LoginSucceeded);
    }

    private void LoginFailed(object data)
    {
        LoginFaild toc = data as LoginFaild;
        Debug.Log("LoginFaild code=" + toc.Code);
    }

    private void LoginSucceeded(object data)
    {
        LoginSuccessfull toc = data as LoginSuccessfull;
        Debug.Log("LoginSuccessfull info=" + toc.PlayerBaseInfo);
    }

    private void LoginRes(object data)
    {
        LoginRes_0X0101 toc = data as LoginRes_0X0101;
        Debug.Log("LoginRes_0X0101 result=" + toc.Result);
    }

    private void STocChat(object data)
    {
        TocChat toc = data as TocChat;
        Debug.Log("STocChat" + toc.Name + toc.Content);
        if (ChatView.Exists)
        {
             string content = toc.Name + ":" + toc.Content;
             Debug.Log(content);
             ChatView.Instance.AddChatItem(content);
        }
    }

    public void CTosChat(string name , string content)
    {
        TosChat tos = new TosChat();
        tos.Name = name;
        tos.Content = content;
        
        SendTos(tos);
    }

    public void LoginReq(string name)
    {
        LoginReq_0X0101 msg = new LoginReq_0X0101();
        msg.AccountName = name; 
        msg.Password = "1111"; 
        SendTos(msg);
        Debug.Log("LoginReq_0X0101 accountName=" + msg.AccountName);
    }
}