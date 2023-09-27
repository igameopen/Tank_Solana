using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using TMPro;
using UnityEngine.UI;
using Common;
using Solana.Unity.SDK;
using Solana.Unity.Rpc.Core.Http;
using Solana.Unity.Wallet;
using System;
using codebase.utility;
using Solana.Unity.SDK.Example;
using Solana.Unity.Rpc.Models;
using Solana.Unity.SDK.Nft;
using System.Globalization;

public class TradingPanel : MonoBehaviour
{
    private TokenAccount _transferTokenAccount = null;
    private Nft _nft = null;

    private PlaneManager planeManager = new PlaneManager();
    private double _ownedSolAmount;
    private const long SolLamports = 1000000000;

    private bool isAgreement;




    private Button m_Button_Close;
    private Button Button_Close
    {
        get
        {
            if (m_Button_Close == null)
                m_Button_Close = transform.Find("Panel/Button_Close").GetComponent<Button>();
            return m_Button_Close;
        }

    }

    private Transform m_TransferPanel;
    private Transform TransferPanel
    {
        get
        {
            if (m_TransferPanel == null)
                m_TransferPanel = transform.Find("Panel/TransferPanel");
            return m_TransferPanel;
        }
    }

    private TMP_InputField m_InputField_Publickey;
    private TMP_InputField InputField_Publickey
    {
        get
        {
            if (m_InputField_Publickey == null)
                m_InputField_Publickey = TransferPanel.Find("Panel/InputField_Publickey").GetComponent<TMP_InputField>();
            return m_InputField_Publickey;
        }
    }

    private TMP_InputField m_InputField_Amount;
    private TMP_InputField InputField_Amount
    {
        get
        {
            if (m_InputField_Amount == null)
                m_InputField_Amount = TransferPanel.Find("Panel/InputField_Amount").GetComponent<TMP_InputField>();
            return m_InputField_Amount;
        }
    }

    private Button m_Button_Paste;
    private Button Button_Paste
    {
        get
        {
            if (m_Button_Paste == null)
                m_Button_Paste = TransferPanel.Find("Panel/Button_Paste").GetComponent<Button>();
            return m_Button_Paste;
        }
    }

    private Button m_Button_Transfer;
    private Button Button_Transfer
    {
        get
        {
            if (m_Button_Transfer == null)
                m_Button_Transfer = TransferPanel.Find("Panel/Button_Transfer").GetComponent<Button>();
            return m_Button_Transfer;
        }
    }

    private Toggle m_Toggle_Agreement;
    private Toggle Toggle_Agreement
    {
        get
        {
            if (m_Toggle_Agreement == null)
                m_Toggle_Agreement = TransferPanel.Find("Panel/AgreementPanel/Toggle_Agreement").GetComponent<Toggle>();
            return m_Toggle_Agreement;
        }
    }

    private Button m_Button_Agreement;
    private Button Button_Agreement
    {
        get
        {
            if (m_Button_Agreement == null)
                m_Button_Agreement = TransferPanel.Find("Panel/AgreementPanel/Button_Agreement").GetComponent<Button>();
            return m_Button_Agreement;
        }
    }




    private Transform m_ReceivePanel;
    private Transform ReceivePanel
    {

        get
        {

            if (m_ReceivePanel == null)
                m_ReceivePanel = transform.Find("Panel/ReceivePanel");
            return m_ReceivePanel;
        }
    }

    private Button m_Button_CopyPublickey;
    private Button Button_CopyPublickey
    {
        get
        {
            if (m_Button_CopyPublickey == null)
                m_Button_CopyPublickey = ReceivePanel.Find("Panel/Button_Publickey").GetComponent<Button>();
            return m_Button_CopyPublickey;
        }
    }

    private TextMeshProUGUI m_Text_Publickey;
    private TextMeshProUGUI Text_Publickey
    {
        get
        {
            if (m_Text_Publickey == null)
                m_Text_Publickey = ReceivePanel.Find("Panel/Button_Publickey/Publickey").GetComponent<TextMeshProUGUI>();
            return m_Text_Publickey;
        }
    }


    private RawImage m_Image_QRCode;
    private RawImage Image_QRCode
    {
        get
        {
            if (m_Image_QRCode == null)
                m_Image_QRCode = ReceivePanel.Find("Panel/QRCodePanel/QRCode").GetComponent<RawImage>();
            return m_Image_QRCode;
        }
    }


    public void Show(TradingType type, TokenAccount TokenAccount = null, Nft nft = null)
    {
        switch (type)
        {
            case TradingType.Transfer:
                TransferPanel.gameObject.SetActive(true);
                ReceivePanel.gameObject.SetActive(false);
                InitTransferPanel(TokenAccount, nft);
                break;

            case TradingType.Receive:
                TransferPanel.gameObject.SetActive(false);
                ReceivePanel.gameObject.SetActive(true);
                InitReceivePanel();
                break;

        }
    }

    async void InitTransferPanel(TokenAccount TokenAccount = null, Nft nft = null)
    {
        if (TokenAccount != null)
        {
            _transferTokenAccount = TokenAccount;
            _ownedSolAmount = double.Parse(_transferTokenAccount.Account.Data.Parsed.Info.TokenAmount.AmountDecimal.ToString(CultureInfo
        .CurrentCulture));
        }
        else
        {
            _ownedSolAmount = await Web3.Instance.WalletBase.GetBalance();
        }
        if (nft != null)
            this._nft = nft;
        Button_Paste.onClick.AddListener(OnButton_Paste);
        Button_Transfer.onClick.AddListener(OnButton_Transfer);

        Toggle_Agreement.onValueChanged.AddListener(OnToggle_Agreement);
        Button_Agreement.onClick.AddListener(OnButton_Agreement);
    }

    void InitReceivePanel()
    {
        Button_CopyPublickey.onClick.AddListener(OnButton_CopyPublickey);
        Texture2D tex = QRGenerator.GenerateQRTexture(Web3.Instance.WalletBase.Account.PublicKey, 256, 256);
        Image_QRCode.texture = tex;

    }

    private void Start()
    {
        Button_Close.onClick.AddListener(OnButton_Close);
    }

    private void OnButton_Close()
    {
        Destroy(gameObject);
    }

    #region 转账

    void OnButton_Paste()
    {
        InputField_Publickey.text = GUIUtility.systemCopyBuffer;
    }

    void OnButton_Transfer()
    {
        if (!isAgreement)
        {
            PrintMSG("Please check to agree to the Risk.");
            return;
        }

        string publickey = InputField_Publickey.text;
        string amount = InputField_Amount.text;
        if (string.IsNullOrEmpty(publickey))
        {
            PrintMSG("Please enter receiver public key");
            return;
        }
        if (string.IsNullOrEmpty(amount))
        {
            PrintMSG("Please input transfer amount");
            return;
        }
        if (float.Parse(InputField_Amount.text) > _ownedSolAmount)
        {
            PrintMSG("Not enough funds for transaction.");
            return;
        }

        if (_nft != null)
        {
            TransferNft();
        }
        else if (_transferTokenAccount == null)
        {
            TransferSol();
        }
        else
        {
            TransferToken();
        }
    }

    private async void TransferSol()
    {
        RequestResult<string> result = await Web3.Instance.WalletBase.Transfer(
            new PublicKey(InputField_Publickey.text),
            Convert.ToUInt64(float.Parse(InputField_Amount.text) * SolLamports));
        HandleResponse(result);
    }

    private async void TransferToken()
    {
        RequestResult<string> result = await Web3.Instance.WalletBase.Transfer(
            new PublicKey(InputField_Publickey.text),
            new PublicKey(_transferTokenAccount.Account.Data.Parsed.Info.Mint),
            Convert.ToUInt64(float.Parse(InputField_Amount.text) * SolLamports));
        HandleResponse(result);
    }
    private async void TransferNft()
    {
        RequestResult<string> result = await Web3.Instance.WalletBase.Transfer(
            new PublicKey(InputField_Publickey.text),
            new PublicKey(_nft.metaplexData.data.mint),
            1);
        HandleResponse(result);
    }



    private void HandleResponse(RequestResult<string> result)
    {
        string error = result.Result == null ? result.Reason : "";
        if (!string.IsNullOrEmpty(error))
        {
            PrintMSG(error);
            return;
        }
        if (result.Result != null)
        {
            Destroy(gameObject);
        }
    }

    private void OnButton_Agreement()
    {
        Application.OpenURL("http://tank.dihub.cn/content/6");
    }

    private void OnToggle_Agreement(bool value)
    {
        isAgreement = value;
    }



    #endregion 转账

    #region 收款

    private void OnButton_CopyPublickey()
    {
        Clipboard.Copy(Text_Publickey.text.Trim());
        PrintMSG("Publickey copied to Clipboard");
    }


    #endregion 收款

    public void PrintMSG(string msg)
    {
        if (!string.IsNullOrEmpty(msg))
        {
            Tools.logText = msg;
            Transform panel = planeManager.Push(new LogPanel()).transform;
            panel.SetParent(transform.parent);
            panel.SetAsLastSibling();

        }
    }

}
