pragma solidity ^0.4.20;

/**
 * 用户向众筹合约转账，记录用户转帐信息
 * 设置有众筹总金额和众筹截至时间
 * 如果截止时间结束后，达到总金额则将众筹款转给众筹项目方，否则退还给用户
 */


contract crowdxy {
	address public beneficiary;//众筹方

//	address public beneficount;//众筹账号

	uint public fundingGoal;   // 募资额度

    uint public amountRaised;   // 众筹总额

    uint public deadline;      // 众筹截止期

    unit public fundnum;       //转账笔数

    mapping(address => uint256) public balanceOf;

    bool fundingGoalReached = false;  // 众筹是否达到目标金额

    bool crowdsaleClosed = false;   //  众筹是否结束

    struct Funder {                 //投资者信息
    	address addr;
    	uint amount;
    	
    }

    Funder[] public funders;

  
    event GoalReached(address addr, uint totalAmountRaised);//达到众筹目标

    event FundTransfer(address addr, uint amount, bool isContribution);//众筹金额有效
    
    
	/**
    * 构造函数 传入投资人地址 众筹目标金额 众筹结束时间
    **/
	function crowdxy (address beneficiaryaddr,
        			  uint fundingGoalInEthers,
                      uint durationInMinutes) {

				beneficiary = beneficiaryaddr;
				fundingGoal = fundingGoalInEthers * 1 ether;
				deadline = now + durationInMinutes * 1 minutes;		
	}	
   /**
     * 无函数名的Fallback函数，
     * 在向合约转账时，这个函数会被调用
     */
    function () payable {
 //       require(!crowdsaleClosed);

        uint amount = msg.value;

        balanceOf[msg.sender] += amount;

        fundnum +=1；

        amountRaised += amount;

        FundTransfer(msg.sender, amount, true);

		funders.push(Funder({

		addr: msg.sender,

		amount: amount

		}));
    }

    /**
    *  
    * 判断是否还在众筹期
    *
    **/
    modifier afterDeadline() { if (now >= deadline) _; }


    /**
     * 判断众筹是否完成融资目标
     *
     */
    function checkGoalReached() afterDeadline {
        if (amountRaised >= fundingGoal) {
            fundingGoalReached = true;
            GoalReached(beneficiary, amountRaised);
        }
        crowdsaleClosed = true;
    }


    /**
     * 完成融资目标时，融资款发送到收款方
     * 未完成融资目标时，执行退款
     *
     */
    function safeWithdrawal() afterDeadline {


        if (!fundingGoalReached) {
            uint amount = balanceOf[msg.sender];
            balanceOf[msg.sender] = 0;
            if (amount > 0) {
                if (msg.sender.send(amount)) {
                    FundTransfer(msg.sender, amount, false);
                } else {
                    balanceOf[msg.sender] = amount;
                }
            }
        }

        if (fundingGoalReached && beneficiary == msg.sender) {
            if (beneficiary.send(amountRaised)) {
                FundTransfer(beneficiary, amountRaised, false);
            } else {
                fundingGoalReached = false;
            }
        }
    }


}
