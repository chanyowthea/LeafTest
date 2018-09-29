using UnityEngine;
using System.Collections;
using Util;
using UnityEngine.UI;

public class ChatView : SingletonMonoBehaviour<ChatView>
{
    public Button btn_chat;
    private InputField inputField;
    private GameObject itemGo;
    [SerializeField] Transform grid;
    private Transform tranf;
	void Start () 
    {
        tranf = transform;
        //btn_chat = tranf.Find("Button").GetComponent<Button>();
        inputField = tranf.Find("InputField").GetComponent<InputField>();
        itemGo = Resources.Load("Prefabs/ChatView/item") as GameObject;
        //UIEventListener.Get(btn_chat.gameObject).onClick = OnClickChatBtn;
	}

    public void OnClickChatBtn(GameObject go)
    {
        //Debug.Log("OnClickChatBtn");
        //AddChatItem(inputField.text);
        //ChatModel.Instance.CTosChat("客户端", inputField.text);
        ChatModel.Instance.LoginReq(inputField.text);
    }

    public void AddChatItem(string content)
    {
        GameObject chatItem = GameObject.Instantiate(itemGo) as GameObject;
        chatItem.transform.SetParent(grid, false);
        chatItem.transform.localPosition = Vector3.zero;
        chatItem.transform.localScale = Vector3.one;
        chatItem.transform.Find("Text").GetComponent<Text>().text = content;
    }

    [SerializeField] InputField _accountNameInput;
    public void OnClickLogin()
    {
        ChatModel.Instance.LoginReq(_accountNameInput.text);
    }
}
