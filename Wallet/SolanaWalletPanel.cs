using Cysharp.Threading.Tasks;
using Solana.Unity.Rpc.Types;
using Solana.Unity.SDK;
using Solana.Unity.SDK.Nft;
using System.Collections;
using System.Collections.Generic;
using System.Linq;
using UnityEngine;
using UnityEngine.UI;
using TMPro;

public class SolanaWalletPanel : MonoBehaviour
{
    public Button button_privateKey;
    public GameObject GetInfoPanel;

    public TextMeshProUGUI SolValue;
    public Button Transfer_Sol;
    public Button Recrive_Btn;


    public GameObject TokenItemPrefab;
    public Transform tokenContainer;

    private List<TokenItem> _instantiatedTokens = new();

    public TradingPanel tradingPanel;

    // Start is called before the first frame update
    void Start()
    {
        button_privateKey.onClick.AddListener(OnPrivateKeyBtn);
        Recrive_Btn.onClick.AddListener(OnRecrive_Btn);
        InitSolItem();
        GetOwnedTokenAccounts();
    }

    private void OnEnable()
    {
        Web3.OnBalanceChange += OnBalanceChange;

    }

    private void OnDisable()
    {
        Web3.OnBalanceChange -= OnBalanceChange;

    }

    // Update is called once per frame
    void Update()
    {

    }

    void OnPrivateKeyBtn()
    {
        Instantiate(GetInfoPanel, transform.root);
    }

    private void OnBalanceChange(double sol)
    {
        SolValue.text = PlayerInfo.instance.solValue;
    }


    void InitSolItem()
    {
        SolValue.text = PlayerInfo.instance.solValue;
        Transfer_Sol.onClick.AddListener(OnTransfer_Sol);
    }

    void OnTransfer_Sol()
    {
        TradingPanel panel = UnityEngine.Object.Instantiate(tradingPanel, transform.root).GetComponent<TradingPanel>();
        panel.Show(Common.TradingType.Transfer);
    }

    void OnRecrive_Btn()
    {
        TradingPanel panel = UnityEngine.Object.Instantiate(tradingPanel, transform.root).GetComponent<TradingPanel>();

        panel.Show(Common.TradingType.Receive);
    }

    private async void GetOwnedTokenAccounts()
    {
        var tokens = await Web3.Wallet.GetTokenAccounts(Commitment.Confirmed);
        if (tokens == null) return;

        if (tokens is { Length: > 0 })
        {
            var tokenAccounts = tokens.OrderByDescending(
                tk => tk.Account.Data.Parsed.Info.TokenAmount.AmountUlong);
            foreach (var item in tokenAccounts)
            {
                if (!(item.Account.Data.Parsed.Info.TokenAmount.AmountUlong > 0)) break;
                if (_instantiatedTokens.All(t => t.TokenAccount.Account.Data.Parsed.Info.Mint != item.Account.Data.Parsed.Info.Mint))
                {
                    var tk = Instantiate(TokenItemPrefab, tokenContainer, true);
                    tk.transform.localScale = Vector3.one;

                    Nft.TryGetNftData(item.Account.Data.Parsed.Info.Mint,
                        Web3.Instance.WalletBase.ActiveRpcClient).AsUniTask().ContinueWith(nft =>
                        {
                            TokenItem tkInstance = tk.GetComponent<TokenItem>();
                            _instantiatedTokens.Add(tkInstance);
                            tk.SetActive(true);
                            if (tkInstance)
                            {
                                tkInstance.InitPanel(item, nft);
                            }
                        }).Forget();
                }
            }
        }



    }

}
